// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attributecache

import (
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

	"github.com/thunder-id/thunderid/internal/system/kmprovider"
)

// Initialize initializes the attribute cache service and returns an instance of AttributeCacheServiceInterface.
func Initialize(
	storeProvider providers.RuntimeStoreProvider,
	crypto kmprovider.RuntimeCryptoProvider,
	encryptionEnabled bool,
) AttributeCacheServiceInterface {
	store := newAttributeCacheStore(storeProvider)
	return newAttributeCacheService(store, crypto, encryptionEnabled)
}
