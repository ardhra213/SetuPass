// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"context"
	"net/http"
	"testing"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/tests/mocks/applicationmock"
	"github.com/thunder-id/thunderid/tests/mocks/design/layoutmock"
	"github.com/thunder-id/thunderid/tests/mocks/design/thememock"
)

// Test Suite
type InitTestSuite struct {
	suite.Suite
}

func TestInitTestSuite(t *testing.T) {
	suite.Run(t, new(InitTestSuite))
}

// Test Initialize returns a non-nil service
func (suite *InitTestSuite) TestInitialize() {
	mux := http.NewServeMux()
	mockTheme := thememock.NewThemeMgtServiceInterfaceMock(suite.T())
	mockLayout := layoutmock.NewLayoutMgtServiceInterfaceMock(suite.T())
	mockApp := applicationmock.NewApplicationServiceInterfaceMock(suite.T())

	service := Initialize(mux, mockTheme, mockLayout, mockApp)

	assert.NotNil(suite.T(), service)
	assert.Implements(suite.T(), (*DesignResolveServiceInterface)(nil), service)
}

// Test registerRoutes does not panic
func (suite *InitTestSuite) TestRegisterRoutes() {
	mux := http.NewServeMux()
	mockService := &mockDesignResolveService{
		resolveDesignFn: func(
			ctx context.Context,
			resolveType providers.DesignResolveType,
			id string,
		) (*providers.DesignResponse, *tidcommon.ServiceError) {
			return nil, nil
		},
	}

	handler := newDesignResolveHandler(mockService)

	// Verify registerRoutes does not panic
	assert.NotPanics(suite.T(), func() {
		registerRoutes(mux, handler)
	})
}
