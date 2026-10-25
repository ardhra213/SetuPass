// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocationcache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/thunder-id/thunderid/internal/system/security"
)

func TestEnforcer_EnsureNotRevoked(t *testing.T) {
	cache := newRevokedCache()
	cache.replace(revokedSnapshot{
		Tokens:   []revokedEntry{{Value: "revoked-jti", ExpiryTime: time.Now().Add(time.Hour)}},
		Families: []revokedEntry{{Value: "revoked-tfid", ExpiryTime: time.Now().Add(time.Hour)}},
		Subjects: []revokedEntry{{Value: "revoked-user", ExpiryTime: time.Now().Add(time.Hour)}},
	})
	e := newEnforcer(cache)

	assert.NoError(t, e.EnsureNotRevoked(context.Background(), security.RevocationIdentity{}),
		"empty ids are a no-op")
	assert.NoError(t, e.EnsureNotRevoked(context.Background(), security.RevocationIdentity{
		JTI: "active-jti", TokenFamilyID: "active-tfid", Subject: "active-user",
	}),
		"a token with a clean jti and family may proceed")
	assert.ErrorIs(t, e.EnsureNotRevoked(context.Background(),
		security.RevocationIdentity{JTI: "revoked-jti"}), errTokenRevoked,
		"a jti on the deny list is rejected")
	assert.ErrorIs(t, e.EnsureNotRevoked(context.Background(), security.RevocationIdentity{
		JTI: "active-jti", TokenFamilyID: "revoked-tfid",
	}), errTokenRevoked,
		"a token whose family is revoked is rejected even with a clean jti")
	assert.ErrorIs(t, e.EnsureNotRevoked(context.Background(), security.RevocationIdentity{
		Subject: "revoked-user",
	}), errTokenRevoked, "a token whose subject is revoked is rejected")
}

func TestNoopEnforcer_AlwaysAllows(t *testing.T) {
	var e EnforcerInterface = noopEnforcer{}
	assert.NoError(t, e.EnsureNotRevoked(context.Background(), security.RevocationIdentity{JTI: "anything"}))
	assert.NoError(t, e.EnsureNotRevoked(context.Background(), security.RevocationIdentity{}))
}
