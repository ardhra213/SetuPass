// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authn

import (
	"sync"
	"testing"

	engineconfig "github.com/thunder-id/thunderid/pkg/thunderidengine/config"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/system/config"
)

type AuthenticationInitTestSuite struct {
	suite.Suite
}

var (
	initRuntimeMutex sync.Mutex
)

func TestAuthenticationInitTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationInitTestSuite))
}

func initializeTestRuntime(root string) error {
	testConfig := &config.Config{
		Server: engineconfig.ServerConfig{
			Hostname: "localhost",
			Port:     8090,
		},
		JWT: engineconfig.JWTConfig{
			Issuer: "test-issuer",
		},
	}
	return config.InitializeServerRuntime(root, testConfig)
}

func (suite *AuthenticationInitTestSuite) SetupSuite() {
	initRuntimeMutex.Lock()
	config.ResetServerRuntime()
	suite.Require().NoError(initializeTestRuntime(suite.T().TempDir()))
}

func (suite *AuthenticationInitTestSuite) TearDownSuite() {
	config.ResetServerRuntime()
	initRuntimeMutex.Unlock()
}
