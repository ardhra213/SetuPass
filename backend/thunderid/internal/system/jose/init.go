// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package jose provides JSON Object Signing and Encryption (JOSE) functionality.
// It includes support for JWS (JSON Web Signature), JWT (JSON Web Token), and JWE (JSON Web Encryption).
package jose

import (
	joseconfig "github.com/thunder-id/thunderid/internal/system/jose/config"
	"github.com/thunder-id/thunderid/internal/system/jose/jwe"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize initializes the JOSE services (JWT and JWE).
func Initialize(
	runtimeProvider providers.RuntimeCryptoProvider, cfg joseconfig.Config,
) (jwt.JWTServiceInterface, jwe.JWEServiceInterface, error) {
	jwtService, err := jwt.Initialize(runtimeProvider, cfg)
	if err != nil {
		return nil, nil, err
	}

	jweService, err := jwe.Initialize(runtimeProvider, cfg)
	if err != nil {
		return nil, nil, err
	}

	return jwtService, jweService, nil
}
