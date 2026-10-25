// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cors

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

// canonicalize returns the canonical form of an HTTP(S) origin for literal-rule
// comparison. The canonical form lowercases the scheme and host, strips a
// trailing dot from the host, Punycode-encodes IDN labels, and wraps IPv6
// hosts in brackets so the resulting string compares equal across all common
// spelling variants. The port is preserved verbatim — operators that want
// both portless and explicit-port forms allowed (e.g. "https://example.com"
// and "https://example.com:443") must list each entry.
//
// This function is used only on the literal-rule path: at compile time on YAML
// literals, and at parse time on the request input via canonicalizeFromURL so
// the matcher hot path does not re-parse. The regex path never invokes this
// function — operator-supplied patterns see the raw header byte for byte.
//
// "null" is not a valid input here; callers must route the null origin through
// the IsNull flag on ParseResult instead.
func canonicalize(raw string) (string, error) {
	if raw == "" {
		return "", fmt.Errorf("%w: empty origin", ErrInvalidOrigin)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidOrigin, err)
	}
	return canonicalizeFromURL(u)
}

// canonicalizeFromURL is the shared core consumed by both Canonicalize and the
// parser's success path so the request hot path avoids a second url.Parse.
func canonicalizeFromURL(u *url.URL) (string, error) {
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case schemeHTTP, schemeHTTPS:
	default:
		return "", fmt.Errorf("%w: unsupported scheme %q", ErrInvalidOrigin, u.Scheme)
	}
	host := strings.ToLower(u.Hostname())
	if host == "" {
		return "", fmt.Errorf("%w: missing host", ErrInvalidOrigin)
	}
	host = strings.TrimSuffix(host, ".")
	if host == "" {
		return "", fmt.Errorf("%w: missing host", ErrInvalidOrigin)
	}
	// Punycode IDN labels so "münchen.example" and "xn--mnchen-3ya.example"
	// compare equal. We accept the lookup result on success and pass through
	// on failure — IDNA-strict rejects valid-but-unusual hosts (underscores,
	// long labels) that browsers happily emit and operators legitimately list.
	if !isIPHost(host) {
		if ascii, err := idna.Lookup.ToASCII(host); err == nil {
			host = ascii
		}
	}
	port := u.Port()
	if strings.Contains(host, ":") {
		// IPv6 host: bracket so the host:port boundary is unambiguous.
		host = "[" + host + "]"
	}
	if port == "" {
		return fmt.Sprintf("%s://%s", scheme, host), nil
	}
	return fmt.Sprintf("%s://%s:%s", scheme, host, port), nil
}

// isIPHost reports whether the given host string is a numeric IP literal
// (IPv4 dotted-quad or IPv6). url.URL.Hostname returns IPv6 addresses without
// brackets, so a colon in the host is a reliable IPv6 signal here.
func isIPHost(host string) bool {
	if strings.Contains(host, ":") {
		return true
	}
	parts := strings.Split(host, ".")
	if len(parts) != 4 {
		return false
	}
	for _, p := range parts {
		if p == "" {
			return false
		}
		for _, c := range p {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}
