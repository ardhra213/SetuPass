// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package jti

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/runtimestore/inmemory"
	"github.com/thunder-id/thunderid/tests/mocks/runtimestoreprovidermock"
)

// JTIStoreTestSuite exercises the jtiStore adapter against a real in-memory runtime store,
// verifying the insert/replay/namespace-isolation semantics.
type JTIStoreTestSuite struct {
	suite.Suite
	store *jtiStore
	ctx   context.Context
}

func TestJTIStoreTestSuite(t *testing.T) {
	suite.Run(t, new(JTIStoreTestSuite))
}

func (suite *JTIStoreTestSuite) SetupTest() {
	suite.store = &jtiStore{storeProvider: inmemory.Initialize("test-deployment")}
	suite.ctx = context.Background()
}

func (suite *JTIStoreTestSuite) TestRecordJTI_Inserted() {
	inserted, err := suite.store.RecordJTI(suite.ctx, "dpop", "jti-1", time.Now().Add(time.Minute))
	suite.Require().NoError(err)
	suite.True(inserted)
}

func (suite *JTIStoreTestSuite) TestRecordJTI_Replay() {
	expiry := time.Now().Add(time.Minute)
	inserted, err := suite.store.RecordJTI(suite.ctx, "dpop", "jti-1", expiry)
	suite.Require().NoError(err)
	suite.True(inserted)

	inserted, err = suite.store.RecordJTI(suite.ctx, "dpop", "jti-1", expiry)
	suite.Require().NoError(err)
	suite.False(inserted, "a repeated jti within the same namespace must be reported as a replay")
}

// TestRecordJTI_NamespaceIsolation locks in the contract that two distinct namespaces can carry
// the same jti without colliding — i.e. namespace participates in the store key.
func (suite *JTIStoreTestSuite) TestRecordJTI_NamespaceIsolation() {
	expiry := time.Now().Add(time.Minute)

	ok1, err := suite.store.RecordJTI(suite.ctx, "dpop", "j", expiry)
	suite.Require().NoError(err)
	suite.True(ok1)

	ok2, err := suite.store.RecordJTI(suite.ctx, "client_assertion", "j", expiry)
	suite.Require().NoError(err)
	suite.True(ok2, "the same jti under a different namespace must not be treated as a replay")
}

func (suite *JTIStoreTestSuite) TestRecordJTI_AlreadyExpired() {
	inserted, err := suite.store.RecordJTI(suite.ctx, "dpop", "jti-1", time.Now().Add(-time.Minute))
	suite.Require().NoError(err)
	suite.True(inserted, "an already-expired jti needs no replay tracking and must not be reported as a replay")
}

func (suite *JTIStoreTestSuite) TestRecordJTI_PutIfNotExistsError() {
	rt := runtimestoreprovidermock.NewRuntimeStoreProviderMock(suite.T())
	rt.EXPECT().PutIfNotExists(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(false, fmt.Errorf("insert failed"))
	store := &jtiStore{storeProvider: rt}

	inserted, err := store.RecordJTI(suite.ctx, "dpop", "jti-1", time.Now().Add(time.Minute))
	suite.Require().Error(err)
	suite.False(inserted)
	suite.Contains(err.Error(), "failed to insert jti")
}
