// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package interceptor provides the interceptor abstraction for cross-cutting flow concerns.
package interceptor

import (
	engineconfig "github.com/thunder-id/thunderid/pkg/thunderidengine/config"
)

// Initialize creates the interceptor registry and registers all built-in interceptors.
func Initialize(
	deps InterceptorDependencies,
	flowConfig engineconfig.FlowConfig,
) (InterceptorRegistryInterface, error) {
	reg := newInterceptorRegistry()
	interceptorNames := flowConfig.Interceptors
	if err := registerInterceptors(deps, reg, interceptorNames); err != nil {
		return nil, err
	}
	// Initialize default interceptorUnits for use in interceptor runner.
	initDefaultInterceptorUnits(deps.FlowFactory)
	return reg, nil
}
