// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package assert

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AssertModelTestSuite struct {
	suite.Suite
}

func TestAssertModelTestSuite(t *testing.T) {
	suite.Run(t, new(AssertModelTestSuite))
}

func (suite *AssertModelTestSuite) TestLevel() {
	testCases := []struct {
		name     string
		level    AssuranceLevel
		expected int
	}{
		{
			name:     "AAL Level 1",
			level:    AALLevel1,
			expected: 1,
		},
		{
			name:     "IAL Level 1",
			level:    IALLevel1,
			expected: 1,
		},
		{
			name:     "AAL Level 2",
			level:    AALLevel2,
			expected: 2,
		},
		{
			name:     "IAL Level 2",
			level:    IALLevel2,
			expected: 2,
		},
		{
			name:     "AAL Level 3",
			level:    AALLevel3,
			expected: 3,
		},
		{
			name:     "IAL Level 3",
			level:    IALLevel3,
			expected: 3,
		},
		{
			name:     "Unknown Level",
			level:    AssuranceLevel("unknown"),
			expected: 0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := tc.level.Level()
			suite.Equal(tc.expected, result)
		})
	}
}
