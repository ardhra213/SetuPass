// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package executor

import (
	engineconfig "github.com/thunder-id/thunderid/pkg/thunderidengine/config"
)

// Initialize creates an executor registry and registers built-in executors.
// When flowConfig.Executors is empty, all built-in executors are registered.
// When non-empty, only the listed executors are registered; flows using other executors
// will fail validation until those executors are included or the list is cleared.
func Initialize(deps ExecutorDependencies, flowConfig engineconfig.FlowConfig) (ExecutorRegistryInterface, error) {
	reg := newExecutorRegistry()
	names := flowConfig.Executors
	if err := registerBuiltInExecutors(reg, deps, names); err != nil {
		return nil, err
	}
	return reg, nil
}
