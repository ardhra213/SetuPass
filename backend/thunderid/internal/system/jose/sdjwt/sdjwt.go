// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package sdjwt

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/thunder-id/thunderid/internal/system/cryptolib"
	"github.com/thunder-id/thunderid/internal/system/jose/jws"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm"
)

// Parse splits a combined-format SD-JWT into the issuer-signed JWT, its
// disclosures, and an optional Key Binding JWT. Disclosure contents are decoded
// and their digests computed under the issuer credential's _sd_alg.
func Parse(combined string) (*Presentation, error) {
	if combined == "" {
		return nil, fmt.Errorf("%w: empty input", ErrInvalidFormat)
	}

	segments := strings.Split(combined, separator)
	if len(segments) < 2 {
		return nil, fmt.Errorf("%w: expected issuer JWT and at least one separator", ErrInvalidFormat)
	}

	issuerJWT := segments[0]
	if issuerJWT == "" {
		return nil, fmt.Errorf("%w: missing issuer JWT", ErrInvalidFormat)
	}

	last := segments[len(segments)-1]
	middle := segments[1 : len(segments)-1]

	var keyBindingJWT string
	if last != "" {
		keyBindingJWT = last
	}

	hashAlg, err := hashAlgFromIssuerJWT(issuerJWT)
	if err != nil {
		return nil, err
	}

	disclosures := make([]Disclosure, 0, len(middle))
	for _, raw := range middle {
		if raw == "" {
			return nil, fmt.Errorf("%w: empty disclosure segment", ErrInvalidFormat)
		}
		d, err := parseDisclosure(raw, hashAlg)
		if err != nil {
			return nil, err
		}
		disclosures = append(disclosures, d)
	}

	return &Presentation{
		IssuerJWT:     issuerJWT,
		Disclosures:   disclosures,
		KeyBindingJWT: keyBindingJWT,
	}, nil
}

// IssuerClaims decodes and returns the issuer-signed JWT payload WITHOUT
// verifying its signature. It is intended for reading the issuer identifier
// (to resolve a trust key) before VerifyIssuerSignature; the returned claims
// MUST NOT be trusted until verification has succeeded.
func (p *Presentation) IssuerClaims() (map[string]interface{}, error) {
	payload, err := decodeJWTPayload(p.IssuerJWT)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	return payload, nil
}

// Verify runs the full SD-JWT verification stack in fail-fast order: issuer
// signature, selective-disclosure resolution, then (when present or required)
// holder binding. On success it returns the resolved credential claims.
func Verify(p *Presentation, opts VerifyOptions) (*VerifiedCredential, error) {
	if err := VerifyIssuerSignature(p, opts.IssuerKey); err != nil {
		return nil, err
	}

	cred, err := ResolveDisclosures(p)
	if err != nil {
		return nil, err
	}

	if p.HasKeyBinding() {
		if err := VerifyKeyBinding(p, cred, opts); err != nil {
			return nil, err
		}
	} else if opts.RequireKeyBinding {
		return nil, ErrMissingKeyBinding
	}

	return cred, nil
}

// VerifyIssuerSignature verifies the issuer-signed JWT against issuerKey.
func VerifyIssuerSignature(p *Presentation, issuerKey crypto.PublicKey) error {
	if issuerKey == nil {
		return fmt.Errorf("%w: no issuer key provided", ErrIssuerSignature)
	}
	if err := verifyJWS(p.IssuerJWT, issuerKey); err != nil {
		return fmt.Errorf("%w: %w", ErrIssuerSignature, err)
	}
	return nil
}

// ResolveDisclosures walks the issuer payload, replaces every referenced digest
// with its disclosed claim, and enforces that all provided disclosures were
// used exactly once.
func ResolveDisclosures(p *Presentation) (*VerifiedCredential, error) {
	payload, err := decodeJWTPayload(p.IssuerJWT)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}

	byDigest := make(map[string]Disclosure, len(p.Disclosures))
	for _, d := range p.Disclosures {
		if _, exists := byDigest[d.Digest]; exists {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateDigest, d.Digest)
		}
		byDigest[d.Digest] = d
	}

	r := &resolver{byDigest: byDigest, used: make(map[string]bool, len(p.Disclosures))}
	resolved, err := r.resolve(payload, "")
	if err != nil {
		return nil, err
	}

	for _, d := range p.Disclosures {
		if !r.used[d.Digest] {
			return nil, fmt.Errorf("%w: %s", ErrUnusedDisclosure, d.Digest)
		}
	}

	claims, ok := resolved.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: issuer payload is not a JSON object", ErrInvalidFormat)
	}

	cred := &VerifiedCredential{Claims: claims, DisclosedPaths: r.disclosed}
	if cnf, ok := claims[claimCnf].(map[string]interface{}); ok {
		if jwk, ok := cnf["jwk"].(map[string]interface{}); ok {
			cred.ConfirmationKey = jwk
		}
	}
	return cred, nil
}

// VerifyKeyBinding verifies the holder binding: the Key Binding JWT is signed by
// the credential's confirmation key and carries the expected typ, aud, nonce,
// iat, and an sd_hash matching the presented issuer JWT and disclosures.
func VerifyKeyBinding(p *Presentation, cred *VerifiedCredential, opts VerifyOptions) error {
	if !p.HasKeyBinding() {
		return ErrMissingKeyBinding
	}
	if cred.ConfirmationKey == nil {
		return ErrMissingConfirmationKey
	}

	holderKey, err := defaultkm.JWKToPublicKey(cred.ConfirmationKey)
	if err != nil {
		return fmt.Errorf("%w: invalid cnf.jwk: %w", ErrMissingConfirmationKey, err)
	}

	header, err := jws.DecodeHeader(p.KeyBindingJWT)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrKeyBindingClaims, err)
	}
	if typ, _ := header["typ"].(string); typ != keyBindingType {
		return fmt.Errorf("%w: unexpected typ %q", ErrKeyBindingClaims, typ)
	}

	if err := verifyJWS(p.KeyBindingJWT, holderKey); err != nil {
		return fmt.Errorf("%w: %w", ErrKeyBindingSignature, err)
	}

	payload, err := decodeJWTPayload(p.KeyBindingJWT)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrKeyBindingClaims, err)
	}

	if opts.ExpectedAudience != "" {
		if aud, _ := payload[claimAud].(string); aud != opts.ExpectedAudience {
			return fmt.Errorf("%w: audience mismatch", ErrKeyBindingClaims)
		}
	}
	if opts.ExpectedNonce != "" {
		if nonce, _ := payload[claimNonce].(string); nonce != opts.ExpectedNonce {
			return fmt.Errorf("%w: nonce mismatch", ErrKeyBindingClaims)
		}
	}

	if err := verifyIat(payload, opts); err != nil {
		return err
	}

	hashAlg, err := hashAlgFromIssuerJWT(p.IssuerJWT)
	if err != nil {
		return err
	}
	expectedSDHash := computeDigest(presentedInput(p), hashAlg)
	if sdHash, _ := payload[claimSDHash].(string); sdHash != expectedSDHash {
		return fmt.Errorf("%w: sd_hash mismatch", ErrKeyBindingClaims)
	}

	return nil
}

// resolver walks a decoded JSON structure resolving SD-JWT digests.
type resolver struct {
	byDigest  map[string]Disclosure
	used      map[string]bool
	disclosed []string
}

func (r *resolver) resolve(node interface{}, path string) (interface{}, error) {
	switch v := node.(type) {
	case map[string]interface{}:
		return r.resolveObject(v, path)
	case []interface{}:
		return r.resolveArray(v, path)
	default:
		return node, nil
	}
}

func (r *resolver) resolveObject(obj map[string]interface{}, path string) (interface{}, error) {
	result := make(map[string]interface{}, len(obj))
	for k, v := range obj {
		if k == claimSD || k == claimSDAlg {
			continue
		}
		resolved, err := r.resolve(v, joinPath(path, k))
		if err != nil {
			return nil, err
		}
		result[k] = resolved
	}

	sdRaw, ok := obj[claimSD]
	if !ok {
		return result, nil
	}
	digests, ok := sdRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: _sd is not an array", ErrInvalidFormat)
	}

	for _, item := range digests {
		digest, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("%w: _sd entry is not a string", ErrInvalidFormat)
		}
		d, ok := r.byDigest[digest]
		if !ok {
			continue // claim not disclosed by the holder
		}
		if r.used[digest] {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateDigest, digest)
		}
		if d.ArrayElement {
			return nil, fmt.Errorf("%w: array disclosure referenced by _sd", ErrInvalidDisclosure)
		}
		if d.ClaimName == claimSD || d.ClaimName == claimArrayElt {
			return nil, fmt.Errorf("%w: reserved claim name %q", ErrClaimCollision, d.ClaimName)
		}
		if _, exists := result[d.ClaimName]; exists {
			return nil, fmt.Errorf("%w: %s", ErrClaimCollision, d.ClaimName)
		}
		r.used[digest] = true
		claimPath := joinPath(path, d.ClaimName)
		r.disclosed = append(r.disclosed, claimPath)
		resolved, err := r.resolve(d.Value, claimPath)
		if err != nil {
			return nil, err
		}
		result[d.ClaimName] = resolved
	}
	return result, nil
}

func (r *resolver) resolveArray(arr []interface{}, path string) (interface{}, error) {
	result := make([]interface{}, 0, len(arr))
	for i, elem := range arr {
		digest, isDigestElem := arrayElementDigest(elem)
		if !isDigestElem {
			resolved, err := r.resolve(elem, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			result = append(result, resolved)
			continue
		}

		d, ok := r.byDigest[digest]
		if !ok {
			continue // array element not disclosed by the holder
		}
		if r.used[digest] {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateDigest, digest)
		}
		if !d.ArrayElement {
			return nil, fmt.Errorf("%w: object disclosure referenced by array element", ErrInvalidDisclosure)
		}
		r.used[digest] = true
		resolved, err := r.resolve(d.Value, fmt.Sprintf("%s[%d]", path, i))
		if err != nil {
			return nil, err
		}
		result = append(result, resolved)
	}
	return result, nil
}

// arrayElementDigest returns the digest of an array element of the form
// {"...": "<digest>"} and whether the element has that shape.
func arrayElementDigest(elem interface{}) (string, bool) {
	m, ok := elem.(map[string]interface{})
	if !ok || len(m) != 1 {
		return "", false
	}
	digest, ok := m[claimArrayElt].(string)
	if !ok {
		return "", false
	}
	return digest, true
}

// parseDisclosure decodes a base64url disclosure into its parts and computes its digest.
func parseDisclosure(raw string, hashAlg string) (Disclosure, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Disclosure{}, fmt.Errorf("%w: base64url decode: %w", ErrInvalidDisclosure, err)
	}

	var arr []interface{}
	if err := json.Unmarshal(decoded, &arr); err != nil {
		return Disclosure{}, fmt.Errorf("%w: not a JSON array: %w", ErrInvalidDisclosure, err)
	}

	d := Disclosure{Raw: raw, Digest: computeDigest(raw, hashAlg)}
	switch len(arr) {
	case 2:
		salt, ok := arr[0].(string)
		if !ok {
			return Disclosure{}, fmt.Errorf("%w: salt is not a string", ErrInvalidDisclosure)
		}
		d.Salt = salt
		d.Value = arr[1]
		d.ArrayElement = true
	case 3:
		salt, ok := arr[0].(string)
		if !ok {
			return Disclosure{}, fmt.Errorf("%w: salt is not a string", ErrInvalidDisclosure)
		}
		name, ok := arr[1].(string)
		if !ok {
			return Disclosure{}, fmt.Errorf("%w: claim name is not a string", ErrInvalidDisclosure)
		}
		d.Salt = salt
		d.ClaimName = name
		d.Value = arr[2]
	default:
		return Disclosure{}, fmt.Errorf("%w: expected 2 or 3 elements, got %d", ErrInvalidDisclosure, len(arr))
	}
	return d, nil
}

// presentedInput is the input to the Key Binding sd_hash: the issuer JWT and all
// disclosures joined by "~" with a trailing "~".
func presentedInput(p *Presentation) string {
	var sb strings.Builder
	sb.WriteString(p.IssuerJWT)
	sb.WriteString(separator)
	for _, d := range p.Disclosures {
		sb.WriteString(d.Raw)
		sb.WriteString(separator)
	}
	return sb.String()
}

// verifyIat validates the Key Binding JWT iat against the reference time.
func verifyIat(payload map[string]interface{}, opts VerifyOptions) error {
	iatRaw, ok := payload[claimIat]
	if !ok {
		return fmt.Errorf("%w: missing iat", ErrKeyBindingClaims)
	}
	iatFloat, ok := iatRaw.(float64)
	if !ok {
		return fmt.Errorf("%w: iat is not a number", ErrKeyBindingClaims)
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	iat := time.Unix(int64(iatFloat), 0)
	if iat.After(now.Add(opts.Leeway)) {
		return fmt.Errorf("%w: iat is in the future", ErrKeyBindingClaims)
	}
	if opts.MaxIATAge > 0 && iat.Before(now.Add(-opts.MaxIATAge)) {
		return fmt.Errorf("%w: iat too old (replay protection)", ErrKeyBindingClaims)
	}
	return nil
}

// hashAlgFromIssuerJWT reads the _sd_alg claim from the issuer JWT, defaulting
// to sha-256 when absent.
func hashAlgFromIssuerJWT(issuerJWT string) (string, error) {
	payload, err := decodeJWTPayload(issuerJWT)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	alg, ok := payload[claimSDAlg].(string)
	if !ok || alg == "" {
		return defaultHashAlg, nil
	}
	switch alg {
	case "sha-256", "sha-384", "sha-512":
		return alg, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedHashAlg, alg)
	}
}

// computeDigest returns base64url(hash(value)) under the named SD-JWT hash algorithm.
func computeDigest(value string, hashAlg string) string {
	var alg cryptolib.HashAlgorithm
	switch hashAlg {
	case "sha-384":
		alg = cryptolib.GenericSHA384
	case "sha-512":
		alg = cryptolib.GenericSHA512
	default:
		alg = cryptolib.GenericSHA256
	}
	sum, err := cryptolib.Hash([]byte(value), alg)
	if err != nil {
		// alg is constrained above to supported values, so this cannot happen.
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(sum)
}

// verifyJWS verifies a compact JWS signature using the given public key.
func verifyJWS(token string, key crypto.PublicKey) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWS format")
	}
	header, err := jws.DecodeHeader(token)
	if err != nil {
		return err
	}
	algStr, _ := header["alg"].(string)
	signAlg, err := mapAlgorithmToSignAlg(algStr)
	if err != nil {
		return err
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}
	signingInput := parts[0] + "." + parts[1]

	return cryptolib.Verify([]byte(signingInput), signature, signAlg, key)
}

// decodeJWTPayload decodes the JSON claims set (the second segment) of a compact JWS.
func decodeJWTPayload(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("base64url decode: %w", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(decoded, &out); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return out, nil
}

// joinPath joins a dotted claim path.
func joinPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}

// mapAlgorithmToSignAlg maps JWS alg header values to internal SignAlgorithm.
func mapAlgorithmToSignAlg(jwsAlg string) (cryptolib.SignAlgorithm, error) {
	signAlg, err := cryptolib.SignAlgorithmFor(cryptolib.Algorithm(defaultkm.JWSAlgorithm(jwsAlg)))
	if err != nil {
		return "", fmt.Errorf("unsupported JWS alg: %s", jwsAlg)
	}
	return signAlg, nil
}
