// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"context"
	"time"
)

// Resolver loads a session from an opaque handle and returns it only when it is live. It does
// not check flow identity or version — that is the SSO-Check node's responsibility.
type Resolver interface {
	// Resolve returns the session referenced by handleID when it is ACTIVE and within its
	// deadlines at now. It returns (nil, nil) for every "no live session" case (absent,
	// ended/revoked, expired), and a non-nil error only on a store failure.
	Resolve(ctx context.Context, handleID string, now time.Time) (*Session, error)
}

type resolver struct {
	store sessionStore
}

// newResolver creates a Resolver backed by the given session store.
func newResolver(store sessionStore) Resolver {
	return &resolver{store: store}
}

// Resolve implements Resolver.
func (r *resolver) Resolve(ctx context.Context, handleID string, now time.Time) (*Session, error) {
	if handleID == "" {
		return nil, nil
	}

	s, err := r.store.GetByHandle(ctx, handleID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	if s.State != StateActive {
		return nil, nil
	}
	if expired(s.IdleExpiresAt, now) || expired(s.AbsoluteExpiresAt, now) {
		return nil, nil
	}

	return s, nil
}

// expired reports whether a deadline is set and has been reached at now. A zero deadline
// means "no deadline" and never expires.
func expired(deadline, now time.Time) bool {
	return !deadline.IsZero() && !now.Before(deadline)
}
