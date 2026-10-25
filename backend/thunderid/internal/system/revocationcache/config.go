// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocationcache

import "time"

// Config holds the Resource Server token-revocation cache settings, mapped by the caller from the
// deployment server security config. It is intentionally decoupled from system/config so this
// package does not depend on the global configuration type.
type Config struct {
	// Enabled turns RS revocation enforcement on. When false, Initialize returns a no-op enforcer
	// and a no-op syncer.
	Enabled bool
	// Source selects the sync source. Currently only "db" is supported; future values may include
	// "endpoint" or "events".
	Source string
	// SyncInterval is how often the cache is refreshed from the source.
	SyncInterval time.Duration
}
