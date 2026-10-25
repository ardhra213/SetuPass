// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package magiclink

const (
	// DefaultExpirySeconds is the default expiry time for magic link tokens in seconds.
	DefaultExpirySeconds = 300

	// tokenAudience is the audience claim for magic link tokens.
	tokenAudience = "magiclink-svc"
)
