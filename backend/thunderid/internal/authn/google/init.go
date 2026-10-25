// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package google

import (
	authnoidc "github.com/thunder-id/thunderid/internal/authn/oidc"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

// Initialize initializes the Google OIDC authentication service.
func Initialize(oidcSvc authnoidc.OIDCAuthnServiceInterface,
	jwtSvc jwt.JWTServiceInterface) GoogleOIDCAuthnServiceInterface {
	return newGoogleOIDCAuthnService(oidcSvc, jwtSvc)
}
