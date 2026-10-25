// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"crypto/tls"
	"testing"

	"github.com/thunder-id/thunderid/internal/system/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UtilsTestSuite defines the test suite for HTTP utils.
type UtilsTestSuite struct {
	suite.Suite
}

// TestUtilsSuite runs the HTTP utils test suite.
func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestGetTLSVersion_TLS12() {
	cfg := config.Config{
		TLS: config.TLSConfig{
			MinVersion: "1.2",
		},
	}

	version := GetTLSVersion(cfg)
	assert.Equal(suite.T(), uint16(tls.VersionTLS12), version)
}

func (suite *UtilsTestSuite) TestGetTLSVersion_TLS13() {
	cfg := config.Config{
		TLS: config.TLSConfig{
			MinVersion: "1.3",
		},
	}

	version := GetTLSVersion(cfg)
	assert.Equal(suite.T(), uint16(tls.VersionTLS13), version)
}

func (suite *UtilsTestSuite) TestGetTLSVersion_DefaultToTLS13() {
	cfg := config.Config{
		TLS: config.TLSConfig{
			MinVersion: "",
		},
	}

	version := GetTLSVersion(cfg)
	assert.Equal(suite.T(), uint16(tls.VersionTLS13), version)
}

func (suite *UtilsTestSuite) TestGetTLSVersion_InvalidVersionDefaultsToTLS13() {
	cfg := config.Config{
		TLS: config.TLSConfig{
			MinVersion: "1.1",
		},
	}

	version := GetTLSVersion(cfg)
	assert.Equal(suite.T(), uint16(tls.VersionTLS13), version)
}

func (suite *UtilsTestSuite) TestGetTLSVersion_UnknownVersionDefaultsToTLS13() {
	cfg := config.Config{
		TLS: config.TLSConfig{
			MinVersion: "invalid",
		},
	}

	version := GetTLSVersion(cfg)
	assert.Equal(suite.T(), uint16(tls.VersionTLS13), version)
}
