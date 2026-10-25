// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"
	"errors"
	"testing"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/tests/mocks/cachemock"
)

const testFlowID = "flow-1"

type GraphCacheTestSuite struct {
	suite.Suite
	cache     GraphCacheInterface
	mockCache *cachemock.CacheInterfaceMock[*graph]
	factory   FlowFactoryInterface
}

func TestGraphCacheTestSuite(t *testing.T) {
	suite.Run(t, new(GraphCacheTestSuite))
}

func (s *GraphCacheTestSuite) SetupTest() {
	s.mockCache = cachemock.NewCacheInterfaceMock[*graph](s.T())
	s.cache = &graphCache{
		cache: s.mockCache,
	}
	s.factory = newFlowFactory()
}

func (s *GraphCacheTestSuite) TestGetSuccess() {
	flowID := testFlowID
	ctx := context.Background()
	g := s.factory.CreateGraph(flowID, providers.FlowTypeAuthentication, 1)

	concreteGraph := g.(*graph)

	s.mockCache.EXPECT().Get(ctx, cache.CacheKey{Key: flowID}).Return(concreteGraph, true)

	result, ok := s.cache.Get(ctx, flowID)

	s.True(ok)
	s.NotNil(result)
	s.Equal(flowID, result.GetID())
}

func (s *GraphCacheTestSuite) TestGetNotFound() {
	flowID := testFlowID
	ctx := context.Background()

	s.mockCache.EXPECT().Get(ctx, cache.CacheKey{Key: flowID}).Return(nil, false)

	result, ok := s.cache.Get(ctx, flowID)

	s.False(ok)
	s.Nil(result)
}

func (s *GraphCacheTestSuite) TestGetEmptyFlowID() {
	result, ok := s.cache.Get(context.Background(), "")

	s.False(ok)
	s.Nil(result)
}

func (s *GraphCacheTestSuite) TestGetNilGraph() {
	flowID := testFlowID
	ctx := context.Background()

	s.mockCache.EXPECT().Get(ctx, cache.CacheKey{Key: flowID}).Return(nil, true)

	result, ok := s.cache.Get(ctx, flowID)

	s.False(ok)
	s.Nil(result)
}

func (s *GraphCacheTestSuite) TestSetSuccess() {
	flowID := testFlowID
	ctx := context.Background()
	g := s.factory.CreateGraph(flowID, providers.FlowTypeAuthentication, 1)

	s.mockCache.EXPECT().Set(ctx, cache.CacheKey{Key: flowID}, g.(*graph)).Return(nil)

	err := s.cache.Set(ctx, flowID, g)

	s.NoError(err)
}

func (s *GraphCacheTestSuite) TestSetCacheError() {
	flowID := testFlowID
	ctx := context.Background()
	g := s.factory.CreateGraph(flowID, providers.FlowTypeAuthentication, 1)

	cacheErr := errors.New("cache error")

	s.mockCache.EXPECT().Set(ctx, cache.CacheKey{Key: flowID}, g.(*graph)).Return(cacheErr)

	err := s.cache.Set(ctx, flowID, g)

	s.Error(err)
	s.Equal(cacheErr, err)
}

func (s *GraphCacheTestSuite) TestSetEmptyFlowID() {
	g := s.factory.CreateGraph(testFlowID, providers.FlowTypeAuthentication, 1)

	err := s.cache.Set(context.Background(), "", g)

	s.Error(err)
	s.Contains(err.Error(), "flowID and graph cannot be empty")
}

func (s *GraphCacheTestSuite) TestSetNilGraph() {
	err := s.cache.Set(context.Background(), testFlowID, nil)

	s.Error(err)
	s.Contains(err.Error(), "flowID and graph cannot be empty")
}

func (s *GraphCacheTestSuite) TestSetEmptyFlowIDAndNilGraph() {
	err := s.cache.Set(context.Background(), "", nil)

	s.Error(err)
	s.Contains(err.Error(), "flowID and graph cannot be empty")
}

func (s *GraphCacheTestSuite) TestSetInvalidGraphType() {
	flowID := testFlowID
	mockGraph := NewGraphInterfaceMock(s.T())

	err := s.cache.Set(context.Background(), flowID, mockGraph)

	s.Error(err)
	s.Contains(err.Error(), "graph must be of concrete type *graph")
}

func (s *GraphCacheTestSuite) TestInvalidateSuccess() {
	flowID := testFlowID
	ctx := context.Background()

	s.mockCache.EXPECT().Delete(ctx, cache.CacheKey{Key: flowID}).Return(nil)

	err := s.cache.Invalidate(ctx, flowID)

	s.NoError(err)
}

func (s *GraphCacheTestSuite) TestInvalidateCacheError() {
	flowID := testFlowID
	ctx := context.Background()
	cacheErr := errors.New("cache error")

	s.mockCache.EXPECT().Delete(ctx, cache.CacheKey{Key: flowID}).Return(cacheErr)

	err := s.cache.Invalidate(ctx, flowID)

	s.Error(err)
	s.Equal(cacheErr, err)
}

func (s *GraphCacheTestSuite) TestInvalidateEmptyFlowID() {
	err := s.cache.Invalidate(context.Background(), "")

	s.NoError(err)
}
