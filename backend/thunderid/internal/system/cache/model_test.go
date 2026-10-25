// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CacheModelTestSuite struct {
	suite.Suite
}

func TestCacheModelSuite(t *testing.T) {
	suite.Run(t, new(CacheModelTestSuite))
}

func (suite *CacheModelTestSuite) TestCacheKeyToString() {
	testCases := []struct {
		name           string
		key            string
		expectedString string
	}{
		{
			name:           "SimpleKey",
			key:            "simple",
			expectedString: "simple",
		},
		{
			name:           "ComplexKey",
			key:            "prefix:complex:key:123",
			expectedString: "prefix:complex:key:123",
		},
		{
			name:           "EmptyKey",
			key:            "",
			expectedString: "",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			cacheKey := CacheKey{Key: tc.key}
			assert.Equal(t, tc.expectedString, cacheKey.ToString())
		})
	}
}

func (suite *CacheModelTestSuite) TestCacheStat() {
	// Test creating and accessing a CacheStat
	stat := CacheStat{
		Enabled:    true,
		Size:       10,
		MaxSize:    100,
		HitCount:   50,
		MissCount:  10,
		HitRate:    0.83,
		EvictCount: 5,
	}

	assert.True(suite.T(), stat.Enabled)
	assert.Equal(suite.T(), 10, stat.Size)
	assert.Equal(suite.T(), 100, stat.MaxSize)
	assert.Equal(suite.T(), int64(50), stat.HitCount)
	assert.Equal(suite.T(), int64(10), stat.MissCount)
	assert.Equal(suite.T(), 0.83, stat.HitRate)
	assert.Equal(suite.T(), int64(5), stat.EvictCount)
}
