// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package security

import (
	"net/http"
)

// AuthenticatorInterface defines the interface for pluggable authentication mechanisms.
// Implementations handle different authentication methods (JWT, API keys, mTLS, etc.).
type AuthenticatorInterface interface {
	// CanHandle determines if this authenticator can process the given request.
	// Returns true if the authenticator recognizes the authentication mechanism in the request.
	CanHandle(r *http.Request) bool

	// Authenticate validates credentials and builds a SecurityContext on success.
	// On success, returns a SecurityContext with authenticated subject information and permissions.
	// On failure, returns an authentication error (401).
	Authenticate(r *http.Request) (*SecurityContext, error)
}
