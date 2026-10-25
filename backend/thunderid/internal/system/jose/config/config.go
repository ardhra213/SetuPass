// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package joseconfig holds JOSE-specific configuration injected at initialization.
package joseconfig

import "time"

// Config holds configuration values required by the JOSE (JWT/JWE) services.
type Config struct {
	Issuer         string
	ValidityPeriod int64
	Audience       string
	PreferredKeyID string
	Leeway         int64
	JWKSCacheTTL   time.Duration
}
