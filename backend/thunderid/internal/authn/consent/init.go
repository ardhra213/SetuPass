// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

// Initialize initializes the consent enforcer service and other related components.
func Initialize(jwtSvc jwt.JWTServiceInterface) ConsentEnforcerService {
	return newConsentEnforcerService(jwtSvc)
}
