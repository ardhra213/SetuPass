// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package security

const (
	// maxPublicPathLength defines the maximum allowed length for a public path.
	// This prevents potential DoS attacks via excessively long paths (even with safe regex).
	maxPublicPathLength = 4096
)
