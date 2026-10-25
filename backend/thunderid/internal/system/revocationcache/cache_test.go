// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocationcache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRevokedCache_ReplaceAndIsRevoked(t *testing.T) {
	c := newRevokedCache()
	future := time.Now().Add(time.Hour)

	assert.False(t, c.isTokenRevoked("jti-1"), "empty cache reports nothing revoked")

	c.replace(revokedSnapshot{
		Tokens: []revokedEntry{
			{Value: "jti-1", ExpiryTime: future},
			{Value: "jti-2", ExpiryTime: future},
		},
		Families: []revokedEntry{{Value: "tfid-1", ExpiryTime: future}},
		Subjects: []revokedEntry{{Value: "user-1", ExpiryTime: future}},
	})

	assert.True(t, c.isTokenRevoked("jti-1"))
	assert.True(t, c.isTokenRevoked("jti-2"))
	assert.False(t, c.isTokenRevoked("jti-3"))
	assert.True(t, c.isTokenFamilyRevoked("tfid-1"))
	assert.False(t, c.isTokenFamilyRevoked("jti-1"), "a jti must not match the family dimension")
	assert.False(t, c.isTokenRevoked("tfid-1"), "a tfid must not match the token dimension")
	assert.True(t, c.isSubjectRevoked("user-1", time.Time{}))
	assert.False(t, c.isSubjectRevoked("jti-1", time.Time{}), "a jti must not match the subject dimension")
}

func TestRevokedCache_ReplaceSwapsSnapshot(t *testing.T) {
	c := newRevokedCache()
	future := time.Now().Add(time.Hour)

	c.replace(revokedSnapshot{Tokens: []revokedEntry{{Value: "old", ExpiryTime: future}}})
	assert.True(t, c.isTokenRevoked("old"))

	c.replace(revokedSnapshot{Tokens: []revokedEntry{{Value: "new", ExpiryTime: future}}})
	assert.False(t, c.isTokenRevoked("old"), "prior entries are dropped on replace")
	assert.True(t, c.isTokenRevoked("new"))
}

func TestRevokedCache_ExpiredEntryNotRevoked(t *testing.T) {
	c := newRevokedCache()
	c.replace(revokedSnapshot{Tokens: []revokedEntry{{Value: "expired", ExpiryTime: time.Now().Add(-time.Second)}}})

	assert.False(t, c.isTokenRevoked("expired"), "an entry past its expiry is treated as not revoked")
}

func TestRevokedCache_SubjectBoundary(t *testing.T) {
	now := time.Now().UTC()
	c := newRevokedCache()
	c.replace(revokedSnapshot{Subjects: []revokedEntry{
		{Value: "bounded-user", ExpiryTime: now.Add(time.Hour), RevokedAt: now, Boundary: true},
		{Value: "deleted-user", ExpiryTime: now.Add(time.Hour), RevokedAt: now},
	}})

	assert.True(t, c.isSubjectRevoked("bounded-user", now.Add(-time.Minute)))
	assert.False(t, c.isSubjectRevoked("bounded-user", now.Add(time.Minute)))
	assert.True(t, c.isSubjectRevoked("deleted-user", now.Add(time.Minute)))
}

func TestRevokedCache_ConcurrentAccess(t *testing.T) {
	c := newRevokedCache()
	future := time.Now().Add(time.Hour)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			c.replace(revokedSnapshot{Tokens: []revokedEntry{{Value: "jti", ExpiryTime: future}}})
		}()
		go func() {
			defer wg.Done()
			_ = c.isTokenRevoked("jti")
		}()
	}
	wg.Wait()

	assert.True(t, c.isTokenRevoked("jti"))
}
