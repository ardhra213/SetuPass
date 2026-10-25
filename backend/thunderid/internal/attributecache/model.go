// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attributecache

// AttributeCache represents a cached attribute entry.
type AttributeCache struct {
	// ID is the unique identifier for the cache entry.
	ID string `json:"id"`

	// Attributes contains the cached attributes.
	Attributes map[string]interface{} `json:"attributes"`

	// TTLSeconds is the time-to-live in seconds for this cache entry.
	TTLSeconds int64 `json:"ttlSeconds"`
}
