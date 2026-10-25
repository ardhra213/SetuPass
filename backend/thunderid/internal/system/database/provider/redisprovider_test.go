// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// RedisProviderTestSuite tests the redisProvider struct methods directly.
//
// Note: initRedisProvider / GetRedisProvider / GetRedisProviderCloser rely on
// a package-level sync.Once and require a live Redis connection. Those paths
// are validated by integration tests against a real Redis server.
type RedisProviderTestSuite struct {
	suite.Suite
}

func TestRedisProviderTestSuite(t *testing.T) {
	suite.Run(t, new(RedisProviderTestSuite))
}

func (suite *RedisProviderTestSuite) TestGetKeyPrefix() {
	p := &redisProvider{keyPrefix: "thunderid"}
	suite.Equal("thunderid", p.GetKeyPrefix())
}

func (suite *RedisProviderTestSuite) TestGetRedisClient_Nil() {
	p := &redisProvider{client: nil}
	suite.Nil(p.GetRedisClient())
}

func (suite *RedisProviderTestSuite) TestClose_NilClient_NoError() {
	// Closing when the client was never initialized should be a no-op.
	p := &redisProvider{client: nil}
	err := p.Close()
	suite.NoError(err)
}
