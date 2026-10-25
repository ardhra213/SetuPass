// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package mgt

import (
	"strings"

	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	declarativeresource "github.com/thunder-id/thunderid/internal/system/declarative_resource"
)

// getI18nStoreMode determines the store mode for translations.
//
// Resolution order:
//  1. If Translation.Store is explicitly configured, use it
//  2. Otherwise, fall back to global DeclarativeResources.Enabled:
//     - If enabled: return "declarative"
//     - If disabled: return "mutable"
//
// Returns normalized store mode: "mutable", "declarative", or "composite"
func getI18nStoreMode(translationConfig config.TranslationConfig) serverconst.StoreMode {
	if translationConfig.Store != "" {
		mode := serverconst.StoreMode(strings.ToLower(strings.TrimSpace(translationConfig.Store)))
		switch mode {
		case serverconst.StoreModeMutable, serverconst.StoreModeDeclarative, serverconst.StoreModeComposite:
			return mode
		}
	}

	if declarativeresource.IsDeclarativeModeEnabled() {
		return serverconst.StoreModeDeclarative
	}

	return serverconst.StoreModeMutable
}

// isDeclarativeModeEnabled reports whether translations run in pure declarative
// (immutable) mode, where all write operations must be rejected. Composite and
// mutable modes allow database-backed writes.
func isDeclarativeModeEnabled() bool {
	return getI18nStoreMode(config.GetServerRuntime().Config.Translation) == serverconst.StoreModeDeclarative
}
