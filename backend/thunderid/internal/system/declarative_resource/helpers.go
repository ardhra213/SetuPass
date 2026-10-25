// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package declarativeresource

import (
	"github.com/thunder-id/thunderid/internal/system/config"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// IsDeclarativeModeEnabled checks if declarative resources are enabled in the configuration.
func IsDeclarativeModeEnabled() bool {
	return config.GetServerRuntime().Config.DeclarativeResources.Enabled
}

// CheckDeclarativeCreate returns an error when declarative read-only mode is active.
func CheckDeclarativeCreate() *tidcommon.ServiceError {
	if IsDeclarativeModeEnabled() {
		return &ErrorDeclarativeResourceCreateOperation
	}
	return nil
}

// CheckDeclarativeUpdate returns an error when declarative read-only mode is active.
func CheckDeclarativeUpdate() *tidcommon.ServiceError {
	if IsDeclarativeModeEnabled() {
		return &ErrorDeclarativeResourceUpdateOperation
	}
	return nil
}

// CheckDeclarativeDelete returns an error when declarative read-only mode is active.
func CheckDeclarativeDelete() *tidcommon.ServiceError {
	if IsDeclarativeModeEnabled() {
		return &ErrorDeclarativeResourceDeleteOperation
	}
	return nil
}
