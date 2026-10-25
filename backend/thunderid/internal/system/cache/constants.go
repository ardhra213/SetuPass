// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cache

// evictionPolicy defines the eviction policy for cache entries.
type evictionPolicy string

const (
	// evictionPolicyLRU represents the Least Recently Used eviction policy.
	evictionPolicyLRU evictionPolicy = "LRU"
	// evictionPolicyLFU represents the Least Frequently Used eviction policy.
	evictionPolicyLFU evictionPolicy = "LFU"
)

// cacheType defines the type of cache being used.
type cacheType string

const (
	// cacheTypeInMemory represents an in-memory cache type.
	cacheTypeInMemory cacheType = "inmemory"
	// cacheTypeRedis represents a Redis-backed cache type.
	cacheTypeRedis cacheType = "redis"
)
