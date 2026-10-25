// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package jwe provides functionalities for handling JSON Web Encryption (JWE).
package jwe

import (
	joseconfig "github.com/thunder-id/thunderid/internal/system/jose/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize initializes the JWE service.
func Initialize(
	cryptoProvider providers.RuntimeCryptoProvider, cfg joseconfig.Config,
) (JWEServiceInterface, error) {
	return newJWEService(cryptoProvider, cfg)
}
