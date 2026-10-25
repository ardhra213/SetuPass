// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package openid4vci

import (
	"errors"
	"net/http"
)

// OpenID4VCI Credential Error Response codes (OpenID4VCI 1.0 §8.3.1) plus the
// OAuth 2.0 codes the issuer endpoints reuse.
const (
	errCodeInvalidCredentialRequest  = "invalid_credential_request" //nolint:gosec
	errCodeUnsupportedCredentialType = "unsupported_credential_type"
	errCodeInvalidProof              = "invalid_proof"
	errCodeInvalidNonce              = "invalid_nonce"
	errCodeInvalidToken              = "invalid_token"
	errCodeInvalidDPoPProof          = "invalid_dpop_proof"
	errCodeServerError               = "server_error"
)

// oid4vciError is the error body a VCI endpoint returns, plus its HTTP status.
// CNonce and CNonceExpiresIn are set on invalid_proof/invalid_nonce responses so
// the wallet can retry with a fresh holder proof (OID4VCI 1.0 §7.2).
type oid4vciError struct {
	Status          int    `json:"-"`
	Code            string `json:"error"`
	Description     string `json:"error_description,omitempty"`
	CNonce          string `json:"c_nonce,omitempty"`
	CNonceExpiresIn int64  `json:"c_nonce_expires_in,omitempty"`
}

// toOID4VCIError maps an internal issuer error to its OpenID4VCI error response.
func toOID4VCIError(err error) oid4vciError {
	switch {
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrUserNotFound):
		return oid4vciError{Status: http.StatusUnauthorized, Code: errCodeInvalidToken,
			Description: "The access token is missing, invalid, or not authorized for the requested credential"}
	case errors.Is(err, ErrInvalidDPoP):
		return oid4vciError{Status: http.StatusUnauthorized, Code: errCodeInvalidDPoPProof,
			Description: "The DPoP proof is missing or does not match the access token"}
	case errors.Is(err, ErrInvalidNonce):
		return oid4vciError{Status: http.StatusBadRequest, Code: errCodeInvalidNonce,
			Description: "The proof carries an unknown or expired c_nonce"}
	case errors.Is(err, ErrInvalidProof):
		return oid4vciError{Status: http.StatusBadRequest, Code: errCodeInvalidProof,
			Description: "The holder proof of possession is missing or invalid"}
	case errors.Is(err, ErrUnsupportedCredential):
		return oid4vciError{Status: http.StatusBadRequest, Code: errCodeUnsupportedCredentialType,
			Description: "No credential configuration is available for the request"}
	case errors.Is(err, ErrInvalidRequest):
		return oid4vciError{Status: http.StatusBadRequest, Code: errCodeInvalidCredentialRequest,
			Description: "The request is missing required parameters or is malformed"}
	default:
		return oid4vciError{Status: http.StatusInternalServerError, Code: errCodeServerError,
			Description: "The request could not be processed"}
	}
}
