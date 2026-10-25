// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package jti

import (
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize returns a JTI replay-cache backend (Redis or relational DB, selected
// by the runtime datasource configuration).
func Initialize(storeProvider providers.RuntimeStoreProvider) JTIStoreInterface {
	return newStore(storeProvider)
}
