// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	authnoauth "github.com/thunder-id/thunderid/internal/authn/oauth"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

// Initialize initializes the OIDC authentication service.
func Initialize(oauthSvc authnoauth.OAuthAuthnServiceInterface,
	jwtSvc jwt.JWTServiceInterface) OIDCAuthnServiceInterface {
	return newOIDCAuthnService(oauthSvc, jwtSvc)
}
