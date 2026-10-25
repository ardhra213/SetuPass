// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package csp builds the deny-first Content-Security-Policy header and the "csp" server-config
// section's value, handler, and runtime resolver.
package csp

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
)

// configSectionCSP duplicates the serverconfig.ConfigNameCSP literal so this package need not import
// serverconfig, which would create an import cycle via the middleware package (CORS does the same).
const configSectionCSP = "csp"

// baselineDirective is a single deny-first directive with its default sources.
type baselineDirective struct {
	name    string
	sources []string
}

// baseline is the deny-first default policy. An effective override replaces a directive's value here
// entirely; an unlisted directive falls back to default-src 'none'.
var baseline = []baselineDirective{
	{"default-src", []string{"'none'"}},
	{"script-src", []string{"'self'"}},
	{"style-src", []string{"'self'"}},
	{"img-src", []string{"'self'"}},
	{"connect-src", []string{"'self'"}},
	{"frame-ancestors", []string{"'none'"}},
	{"base-uri", []string{"'none'"}},
	{"form-action", []string{"'self'"}},
}

// directiveNamePattern validates a CSP directive name, guarding against header injection via a crafted
// name rather than restricting which directives may be configured.
var directiveNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// invalidValueChars are disallowed in a directive source or report_uri, so a configured value cannot
// inject an extra directive into the emitted header.
const invalidValueChars = " \t\r\n;,"

// PathPolicy overrides directives for requests whose path starts with Location, applied over the
// policy's default directives and the baseline.
type PathPolicy struct {
	Location   string              `json:"location" yaml:"location"`
	Directives map[string][]string `json:"directives" yaml:"directives"`
}

// PolicyConfig is the value of the "csp" server-config section. For a request, the longest-matching
// Paths entry's directives override Directives (the default), which override the baseline. ReportOnly
// is a pointer so an unset value defaults to report-only rather than silently enforcing.
type PolicyConfig struct {
	ReportOnly *bool               `json:"reportOnly,omitempty" yaml:"reportOnly,omitempty"`
	ReportURI  string              `json:"reportUri,omitempty"  yaml:"reportUri,omitempty"`
	Directives map[string][]string `json:"directives,omitempty" yaml:"directives,omitempty"`
	Paths      []PathPolicy        `json:"paths,omitempty"      yaml:"paths,omitempty"`
}

// EffectiveReportOnly defaults to true when unset, so a policy is never enforced until explicitly
// switched.
func (c PolicyConfig) EffectiveReportOnly() bool {
	return c.ReportOnly == nil || *c.ReportOnly
}

// HeaderName returns the report-only or enforcing header name, per EffectiveReportOnly.
func (c PolicyConfig) HeaderName() string {
	if c.EffectiveReportOnly() {
		return serverconst.ContentSecurityPolicyReportOnlyHeaderName
	}
	return serverconst.ContentSecurityPolicyHeaderName
}

// effectiveDirectives returns the directive overrides for a request path: Directives, with the
// longest-matching Paths entry applied on top, per directive.
func (c PolicyConfig) effectiveDirectives(path string) map[string][]string {
	eff := make(map[string][]string, len(c.Directives))
	for name, sources := range c.Directives {
		eff[name] = sources
	}
	bestLen := -1
	var best map[string][]string
	for i := range c.Paths {
		p := c.Paths[i]
		if strings.HasPrefix(path, p.Location) && len(p.Location) > bestLen {
			bestLen = len(p.Location)
			best = p.Directives
		}
	}
	for name, sources := range best {
		eff[name] = sources
	}
	return eff
}

// nonceDirectives are the directives the per-request nonce is appended to. style-src-attr (inline
// style="" attributes) is deliberately excluded: attributes can't carry a nonce per spec, and it needs
// 'unsafe-inline' for React/MUI's style={{}} usage, which a nonce elsewhere wouldn't affect since it's
// a separate directive.
var nonceDirectives = map[string]struct{}{"script-src": {}, "style-src-elem": {}}

// sourcesWithNonce appends a 'nonce-<value>' token to sources for a nonce directive, without mutating
// the caller's slice (a config value shared across concurrent requests). Leaves sources untouched if
// 'unsafe-inline' is already there: per spec, a nonce present in a directive makes browsers ignore
// 'unsafe-inline' in that same directive, so appending one would defeat an operator's relaxation
// instead of tightening it (see csp.yaml's Console entry for why that relaxation exists).
func sourcesWithNonce(name string, sources []string, nonce string) []string {
	if _, ok := nonceDirectives[name]; !ok || nonce == "" {
		return sources
	}
	for _, source := range sources {
		if source == "'unsafe-inline'" {
			return sources
		}
	}
	withNonce := make([]string, len(sources)+1)
	copy(withNonce, sources)
	withNonce[len(sources)] = "'nonce-" + nonce + "'"
	return withNonce
}

// HeaderValue renders the effective policy for a request path, followed by a report-uri directive when
// a report endpoint is set. nonce, generated fresh per request by the middleware, is appended to each
// directive in nonceDirectives.
func (c PolicyConfig) HeaderValue(path string, nonce string) string {
	overrides := c.effectiveDirectives(path)
	baselineNames := make(map[string]struct{}, len(baseline))
	directives := make([]string, 0, len(baseline)+len(overrides)+1)

	for _, d := range baseline {
		baselineNames[d.name] = struct{}{}
		sources := d.sources
		if override, ok := overrides[d.name]; ok {
			sources = override
		}
		sources = sourcesWithNonce(d.name, sources, nonce)
		directives = append(directives, d.name+" "+strings.Join(sources, " "))
	}

	// Effective directives outside the baseline, in a stable order.
	extraNames := make([]string, 0, len(overrides))
	for name := range overrides {
		if _, isBaseline := baselineNames[name]; !isBaseline {
			extraNames = append(extraNames, name)
		}
	}
	sort.Strings(extraNames)
	for _, name := range extraNames {
		sources := sourcesWithNonce(name, overrides[name], nonce)
		directives = append(directives, name+" "+strings.Join(sources, " "))
	}

	if c.ReportURI != "" {
		directives = append(directives, "report-uri "+c.ReportURI)
	}
	return strings.Join(directives, "; ")
}

// Validate reports whether the configured policy is acceptable. No directive is locked; a report_uri
// must be a relative path or an http/https URL; each path policy needs a rooted prefix and a directive.
func (c PolicyConfig) Validate() error {
	if err := validateDirectives(c.Directives); err != nil {
		return err
	}
	for _, p := range c.Paths {
		if !strings.HasPrefix(p.Location, "/") {
			return fmt.Errorf("csp: location %q must start with %q", p.Location, "/")
		}
		if len(p.Directives) == 0 {
			return fmt.Errorf("csp: path %q has no directives", p.Location)
		}
		if err := validateDirectives(p.Directives); err != nil {
			return fmt.Errorf("csp: path %q: %w", p.Location, err)
		}
	}

	if c.ReportURI != "" {
		if strings.ContainsAny(c.ReportURI, invalidValueChars) {
			return fmt.Errorf("csp: report_uri has an invalid character")
		}
		parsed, err := url.Parse(c.ReportURI)
		if err != nil {
			return fmt.Errorf("csp: report_uri is not a valid URL: %w", err)
		}
		if parsed.Scheme != "" && parsed.Scheme != "http" && parsed.Scheme != "https" {
			return fmt.Errorf("csp: report_uri must use http or https scheme (got %q)", parsed.Scheme)
		}
	}
	return nil
}

// validateDirectives rejects invalid directive names and sources containing a separator or whitespace
// character, so a configured value cannot inject additional directives into the header.
func validateDirectives(directives map[string][]string) error {
	names := make([]string, 0, len(directives))
	for name := range directives {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !directiveNamePattern.MatchString(name) {
			return fmt.Errorf("csp: %q is not a valid directive name", name)
		}
		if len(directives[name]) == 0 {
			return fmt.Errorf("csp: directive %q has no sources", name)
		}
		for _, source := range directives[name] {
			if strings.TrimSpace(source) == "" {
				return fmt.Errorf("csp: directive %q has an empty source", name)
			}
			if strings.ContainsAny(source, invalidValueChars) {
				return fmt.Errorf("csp: directive %q has a source with an invalid character", name)
			}
		}
	}
	return nil
}
