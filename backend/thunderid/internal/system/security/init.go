// Copyright 2025-2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package security

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

// Initialize creates and returns the security middleware with necessary authenticators. The
// revocationEnforcer is consulted after authentication to reject revoked tokens.
func Initialize(jwtService jwt.JWTServiceInterface, revocationEnforcer RevocationEnforcerInterface,
) (func(http.Handler) http.Handler, error) {
	jwtAuthenticator := newJWTAuthenticator(jwtService)
	securityService, err := newSecurityService(
		[]AuthenticatorInterface{jwtAuthenticator}, revocationEnforcer, publicPaths, apiPermissionEntries)
	if err != nil {
		return nil, err
	}
	return middleware(securityService)
}
