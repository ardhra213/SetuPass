// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package thememgt

import (
	"context"
	"strings"

	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	declarativeresource "github.com/thunder-id/thunderid/internal/system/declarative_resource"
	"github.com/thunder-id/thunderid/internal/system/log"
)

// getThemeStoreMode determines the store mode for themes.
//
// Resolution order:
//  1. If Theme.Store is explicitly configured, use it
//  2. Otherwise, fall back to global DeclarativeResources.Enabled:
//     - If enabled: return "declarative"
//     - If disabled: return "mutable"
//
// Returns normalized store mode: "mutable", "declarative", or "composite"
func getThemeStoreMode() serverconst.StoreMode {
	cfg := config.GetServerRuntime().Config
	// Check if service-level configuration is explicitly set
	if cfg.Theme.Store != "" {
		mode := serverconst.StoreMode(strings.ToLower(strings.TrimSpace(cfg.Theme.Store)))
		// Validate and normalize
		switch mode {
		case serverconst.StoreModeMutable, serverconst.StoreModeDeclarative, serverconst.StoreModeComposite:
			return mode
		default:
			// Warn about unrecognized value and fall back to default
			logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ThemeMgtConfig"))
			// Store mode is resolved at server startup, outside any request,
			// so there is no request context (or trace ID) to propagate.
			logger.Warn(context.Background(), "Unrecognized theme store configuration value",
				log.String("raw_value", cfg.Theme.Store),
				log.String("normalized_value", string(mode)),
				log.String("fallback", "global declarative_resources setting"))
		}
	}

	// Fall back to global declarative resources setting
	if declarativeresource.IsDeclarativeModeEnabled() {
		return serverconst.StoreModeDeclarative
	}

	return serverconst.StoreModeMutable
}

// isDeclarativeModeEnabled checks if immutable-only store mode is enabled for themes.
func isDeclarativeModeEnabled() bool {
	return getThemeStoreMode() == serverconst.StoreModeDeclarative
}
