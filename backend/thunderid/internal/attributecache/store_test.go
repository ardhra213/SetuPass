// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attributecache

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/runtimestore/inmemory"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// AttributeCacheStoreTestSuite exercises the attributeCacheStore adapter against a real in-memory
// runtime store, verifying the namespace/key/TTL round-trip and the not-found semantics. The store
// persists opaque bytes; serialization and encryption live in the service layer.
type AttributeCacheStoreTestSuite struct {
	suite.Suite
	store    attributeCacheStoreInterface
	ctx      context.Context
	testID   string
	testData []byte
}

func TestAttributeCacheStoreSuite(t *testing.T) {
	suite.Run(t, new(AttributeCacheStoreTestSuite))
}

func (suite *AttributeCacheStoreTestSuite) SetupTest() {
	suite.store = newAttributeCacheStore(inmemory.Initialize("test-deployment"))
	suite.ctx = context.Background()
	suite.testID = "test-cache-id"
	suite.testData = []byte(`{"key":"value"}`)
}

// Tests for CreateAttributeCache

func (suite *AttributeCacheStoreTestSuite) TestCreateAttributeCache_Success() {
	err := suite.store.CreateAttributeCache(suite.ctx, suite.testID, suite.testData, 3600)
	suite.Require().NoError(err)

	got, err := suite.store.GetAttributeCache(suite.ctx, suite.testID)
	suite.Require().NoError(err)
	suite.Equal(suite.testData, got)
}

// Tests for GetAttributeCache

func (suite *AttributeCacheStoreTestSuite) TestGetAttributeCache_NotFound() {
	got, err := suite.store.GetAttributeCache(suite.ctx, "missing")

	suite.ErrorIs(err, errAttributeCacheNotFound)
	suite.Nil(got)
}

// Tests for ExtendAttributeCacheTTL

func (suite *AttributeCacheStoreTestSuite) TestExtendAttributeCacheTTL_Success() {
	suite.Require().NoError(suite.store.CreateAttributeCache(suite.ctx, suite.testID, suite.testData, 3600))

	err := suite.store.ExtendAttributeCacheTTL(suite.ctx, suite.testID, 7200)
	suite.Require().NoError(err)

	got, err := suite.store.GetAttributeCache(suite.ctx, suite.testID)
	suite.Require().NoError(err)
	suite.Equal(suite.testData, got)
}

func (suite *AttributeCacheStoreTestSuite) TestExtendAttributeCacheTTL_NotFound() {
	err := suite.store.ExtendAttributeCacheTTL(suite.ctx, "missing", 7200)

	suite.Error(err)
}

// Tests for DeleteAttributeCache

func (suite *AttributeCacheStoreTestSuite) TestDeleteAttributeCache_Success() {
	suite.Require().NoError(suite.store.CreateAttributeCache(suite.ctx, suite.testID, suite.testData, 3600))
	suite.Require().NoError(suite.store.DeleteAttributeCache(suite.ctx, suite.testID))

	_, err := suite.store.GetAttributeCache(suite.ctx, suite.testID)
	suite.ErrorIs(err, errAttributeCacheNotFound)
}

func (suite *AttributeCacheStoreTestSuite) TestDeleteAttributeCache_NotFound() {
	err := suite.store.DeleteAttributeCache(suite.ctx, "missing")

	suite.NoError(err)
}

// erroringRuntimeStore wraps a RuntimeStoreProvider and forces Get to fail, so that
// GetAttributeCache's error-propagation path (as opposed to its not-found path) can be exercised.
type erroringRuntimeStore struct {
	providers.RuntimeStoreProvider
}

func (e *erroringRuntimeStore) Get(_ context.Context, _ providers.RuntimeStoreNamespace,
	_ string) ([]byte, error) {
	return nil, errors.New("store unavailable")
}

func (suite *AttributeCacheStoreTestSuite) TestGetAttributeCache_StoreError() {
	s := newAttributeCacheStore(&erroringRuntimeStore{RuntimeStoreProvider: inmemory.Initialize("test-deployment")})

	_, err := s.GetAttributeCache(suite.ctx, suite.testID)

	suite.Error(err)
	suite.Contains(err.Error(), "failed to get attribute cache")
	suite.False(errors.Is(err, errAttributeCacheNotFound))
}
