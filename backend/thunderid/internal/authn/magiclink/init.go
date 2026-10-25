// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package magiclink

import (
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

// Initialize initializes the Magic Link authentication service.
func Initialize(
	jwtSvc jwt.JWTServiceInterface,
) MagicLinkAuthnServiceInterface {
	return newMagicLinkAuthnService(jwtSvc)
}
