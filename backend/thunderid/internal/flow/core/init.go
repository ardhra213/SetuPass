// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package core provides the core structs for flow management and execution.
package core

import "github.com/thunder-id/thunderid/internal/system/cache"

// Initialize initializes the core flow package
func Initialize(cacheManager cache.CacheManagerInterface) (FlowFactoryInterface, GraphCacheInterface) {
	flowFactory := newFlowFactory()
	graphCacheInst := cache.GetInMemoryCache[*graph](cacheManager, "FlowGraphCache")
	graphCache := newGraphCache(graphCacheInst)
	return flowFactory, graphCache
}
