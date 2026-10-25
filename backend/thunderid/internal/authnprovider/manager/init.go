// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package manager

import (
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize creates a new AuthnProviderManager. defaultProvider is required.
func Initialize(defaultProvider providers.AuthnProviderInterface,
	others map[string]providers.CustomAuthnProvider) (providers.AuthnProviderManager, error) {
	return newAuthnProviderManager(defaultProvider, others)
}
