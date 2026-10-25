// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

type ModelTestSuite struct {
	suite.Suite
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

func (s *ModelTestSuite) TestEntityCategoryString() {
	s.Equal("user", providers.EntityCategoryUser.String())
	s.Equal("app", providers.EntityCategoryApp.String())
}

func (s *ModelTestSuite) TestEntityStateString() {
	s.Equal("ACTIVE", providers.EntityStateActive.String())
}
