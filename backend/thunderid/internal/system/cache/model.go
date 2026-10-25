// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"time"
)

// CacheKey represents a key for the cache.
type CacheKey struct {
	Key string
}

// ToString returns the string representation of the CacheKey.
func (key CacheKey) ToString() string {
	return key.Key
}

// CacheEntry represents a cache entry.
type CacheEntry[T any] struct {
	Value      T
	ExpiryTime time.Time
}

// CacheStat represents cache statistics.
type CacheStat struct {
	Enabled    bool
	Size       int
	MaxSize    int
	HitCount   int64
	MissCount  int64
	HitRate    float64
	EvictCount int64
}
