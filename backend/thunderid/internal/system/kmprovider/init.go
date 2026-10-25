// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package kmprovider provides the key manager provider abstraction and initialization.
package kmprovider

import (
	"github.com/thunder-id/thunderid/internal/system/kmprovider/common"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm/pki"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// RuntimeCryptoProvider is a type alias for convenience.
type RuntimeCryptoProvider = providers.RuntimeCryptoProvider

// ConfigCryptoProvider is a type alias for convenience.
type ConfigCryptoProvider = common.ConfigCryptoProvider

// Initialize initializes and returns both RuntimeCryptoProvider and ConfigCryptoProvider.
// The pkiService is injected as a dependency.
// Currently hardcoded to use the default KM provider, but structured to support
// provider selection based on server configuration in the future.
func Initialize(
	pkiService pki.PKIServiceInterface,
) (providers.RuntimeCryptoProvider, common.ConfigCryptoProvider, error) {
	return defaultkm.Initialize(pkiService)
}
