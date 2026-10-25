// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package jwks

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type JWKSHandlerTestSuite struct {
	suite.Suite
	mockService *JWKSServiceInterfaceMock
	handler     *jwksHandler
}

func TestJWKSHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(JWKSHandlerTestSuite))
}

func (s *JWKSHandlerTestSuite) SetupTest() {
	s.mockService = NewJWKSServiceInterfaceMock(s.T())
	s.handler = newJWKSHandler(s.mockService)
}

func (s *JWKSHandlerTestSuite) TestNewJWKSHandler() {
	handler := newJWKSHandler(s.mockService)
	assert.NotNil(s.T(), handler)
	assert.NotNil(s.T(), handler.jwksService)
}

func (s *JWKSHandlerTestSuite) TestHandleJWKSRequest_Success() {
	req := httptest.NewRequest(http.MethodGet, "/oauth2/jwks", nil)
	rr := httptest.NewRecorder()

	jwksResponse := &JWKSResponse{
		Keys: []JWKS{
			{
				Kid: "test-kid",
				Kty: "RSA",
				Use: "sig",
				Alg: "RS256",
				N:   "test-n",
				E:   "AQAB",
			},
		},
	}
	s.mockService.On("GetJWKS", mock.Anything).Return(jwksResponse, nil)

	s.handler.HandleJWKSRequest(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), "application/json", rr.Header().Get("Content-Type"))

	var response JWKSResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), response.Keys, 1)
	assert.Equal(s.T(), "test-kid", response.Keys[0].Kid)
	s.mockService.AssertExpectations(s.T())
}

func (s *JWKSHandlerTestSuite) TestHandleJWKSRequest_ClientError() {
	req := httptest.NewRequest(http.MethodGet, "/oauth2/jwks", nil)
	rr := httptest.NewRecorder()

	svcErr := &tidcommon.ServiceError{
		Type:             tidcommon.ClientErrorType,
		Code:             "invalid_request",
		Error:            tidcommon.I18nMessage{Key: "error.test.invalid_request", DefaultValue: "invalid_request"},
		ErrorDescription: tidcommon.I18nMessage{Key: "error.test.invalid_request", DefaultValue: "Invalid request"},
	}
	s.mockService.On("GetJWKS", mock.Anything).Return(nil, svcErr)

	s.handler.HandleJWKSRequest(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Equal(s.T(), "application/json", rr.Header().Get("Content-Type"))
	s.mockService.AssertExpectations(s.T())
}

func (s *JWKSHandlerTestSuite) TestHandleJWKSRequest_ServiceError() {
	req := httptest.NewRequest(http.MethodGet, "/oauth2/jwks", nil)
	rr := httptest.NewRecorder()

	svcErr := tidcommon.CustomServiceError(tidcommon.InternalServerError, tidcommon.I18nMessage{
		Key:          "error.test.failed_get_jwks",
		DefaultValue: "Failed to get JWKS",
	})
	s.mockService.On("GetJWKS", mock.Anything).Return(nil, svcErr)

	s.handler.HandleJWKSRequest(rr, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), "application/json", rr.Header().Get("Content-Type"))
	assert.Contains(s.T(), rr.Body.String(), svcErr.Code)
	s.mockService.AssertExpectations(s.T())
}
