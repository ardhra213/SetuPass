// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package openid4vci

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/system/config"
)

type InitTestSuite struct {
	suite.Suite
}

func TestInitTestSuite(t *testing.T) {
	suite.Run(t, new(InitTestSuite))
}

func (s *InitTestSuite) TestRegisterRoutes() {
	svc := NewOpenID4VCIServiceInterfaceMock(s.T())
	svc.EXPECT().GetMetadata(mock.Anything).Return(map[string]interface{}{"credential_issuer": "https://i"}).Maybe()
	svc.EXPECT().GenerateNonce(mock.Anything).Return("nonce", nil).Maybe()

	mux := http.NewServeMux()
	registerRoutes(mux, newOpenID4VCIHandler(svc, nil, "https://i/credential", time.Minute))

	cases := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, metadataPath, http.StatusOK},
		{http.MethodPost, noncePath, http.StatusOK},
		{http.MethodOptions, metadataPath, http.StatusNoContent},
		{http.MethodOptions, credentialOfferPath + "/abc", http.StatusNoContent},
	}
	for _, c := range cases {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(c.method, c.path, nil))
		s.Equal(c.status, rr.Code, "%s %s", c.method, c.path)
	}
}

// Initialize disables the issuer engine (nil service) when no signing key is configured.
func (s *InitTestSuite) TestInitializeDisabledWithoutSigningKey() {
	config.ResetServerRuntime()
	s.Require().NoError(config.InitializeServerRuntime("", &config.Config{}))
	defer config.ResetServerRuntime()

	svc, err := Initialize(http.NewServeMux(), nil, nil, nil, nil, nil, nil, nil)
	s.Require().NoError(err)
	s.Nil(svc)
}
