// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package serverconfig

import (
	"fmt"
	"strings"

	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	declarativeresource "github.com/thunder-id/thunderid/internal/system/declarative_resource"
)

// getServerConfigStoreMode determines the store mode for server config.
//
// Resolution order:
//  1. If ServerConfig.Store is explicitly configured, use it; an unrecognized value is an error so a
//     misconfiguration fails fast at startup rather than silently changing store behavior.
//  2. Otherwise, fall back to the global declarative-resources setting: enabled → composite (declarative
//     defaults plus runtime overrides), disabled → mutable.
//
// The writable layer is present in both fallback modes (mutable and composite), so PUT is always allowed
// there; an explicit "declarative" makes the section read-only.
func getServerConfigStoreMode() (serverconst.StoreMode, error) {
	if store := config.GetServerRuntime().Config.ServerConfig.Store; store != "" {
		switch serverconst.StoreMode(strings.ToLower(strings.TrimSpace(store))) {
		case serverconst.StoreModeMutable:
			return serverconst.StoreModeMutable, nil
		case serverconst.StoreModeDeclarative:
			return serverconst.StoreModeDeclarative, nil
		case serverconst.StoreModeComposite:
			return serverconst.StoreModeComposite, nil
		default:
			return "", fmt.Errorf("serverconfig: invalid server_config.store %q", store)
		}
	}
	if declarativeresource.IsDeclarativeModeEnabled() {
		return serverconst.StoreModeComposite, nil
	}
	return serverconst.StoreModeMutable, nil
}
