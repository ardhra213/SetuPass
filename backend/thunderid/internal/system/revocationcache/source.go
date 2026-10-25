// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocationcache

import "context"

// syncSource supplies the current deny-list snapshot to the cache. It is the pluggable seam that lets
// the Resource Server sync from the runtime persistent DB today and from another DB, endpoint, or event stream
// in future without changing the cache, enforcer, or syncer.
type syncSource interface {
	// Snapshot returns all currently-revoked, non-expired entries for this deployment: the revoked
	// single-token jtis and the revoked token-family ids.
	Snapshot(ctx context.Context) (revokedSnapshot, error)
}
