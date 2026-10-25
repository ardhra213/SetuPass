// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package thememgt

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
)

// ConfigTestSuite tests theme configuration functions.
type ConfigTestSuite struct {
	suite.Suite
}

// SetupSuite sets up the test suite once.
func (suite *ConfigTestSuite) SetupSuite() {
	// Reset server runtime to ensure clean state
	config.ResetServerRuntime()
	// Initialize runtime once for all tests in the suite
	testConfig := &config.Config{}
	err := config.InitializeServerRuntime("/tmp/test", testConfig)
	if err != nil {
		suite.Fail("Failed to initialize runtime", err)
	}
}

// TearDownSuite cleans up after the test suite.
func (suite *ConfigTestSuite) TearDownSuite() {
	// Reset server runtime to avoid state leakage to other test suites
	config.ResetServerRuntime()
}

// SetupTest sets up the test environment before each test.
func (suite *ConfigTestSuite) SetupTest() {
	// Reset config before each test
	runtime := config.GetServerRuntime()
	runtime.Config.Theme.Store = ""
	runtime.Config.DeclarativeResources.Enabled = false
}

// TestGetThemeStoreMode tests the store mode resolution logic.
func (suite *ConfigTestSuite) TestGetThemeStoreMode() {
	testCases := []struct {
		name         string
		themeStore   string
		declEnabled  bool
		expectedMode serverconst.StoreMode
		description  string
	}{
		{
			name:         "explicit mutable mode",
			themeStore:   "mutable",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store is explicitly set to 'mutable'",
		},
		{
			name:         "explicit declarative mode",
			themeStore:   "declarative",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeDeclarative,
			description:  "when theme.store is explicitly set to 'declarative'",
		},
		{
			name:         "explicit composite mode",
			themeStore:   "composite",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeComposite,
			description:  "when theme.store is explicitly set to 'composite'",
		},
		{
			name:         "case insensitive mutable",
			themeStore:   "MUTABLE",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store is case-insensitive",
		},
		{
			name:         "case insensitive declarative",
			themeStore:   "DECLARATIVE",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeDeclarative,
			description:  "when theme.store is case-insensitive",
		},
		{
			name:         "case insensitive composite",
			themeStore:   "Composite",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeComposite,
			description:  "when theme.store is case-insensitive",
		},
		{
			name:         "trimmed whitespace mutable",
			themeStore:   "  mutable  ",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store has leading/trailing whitespace",
		},
		{
			name:         "fallback to declarative when enabled",
			themeStore:   "",
			declEnabled:  true,
			expectedMode: serverconst.StoreModeDeclarative,
			description:  "when theme.store is empty and declarative_resources.enabled is true",
		},
		{
			name:         "fallback to mutable when disabled",
			themeStore:   "",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store is empty and declarative_resources.enabled is false",
		},
		{
			name:         "explicit config takes precedence over global setting",
			themeStore:   "mutable",
			declEnabled:  true,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store is set, it takes precedence over global setting",
		},
		{
			name:         "invalid mode fallback to mutable",
			themeStore:   "invalid_mode",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store has invalid value, fallback to global setting (mutable)",
		},
		{
			name:         "invalid mode fallback to declarative",
			themeStore:   "unknown_value",
			declEnabled:  true,
			expectedMode: serverconst.StoreModeDeclarative,
			description:  "when theme.store has invalid value, fallback to global declarative setting",
		},
		{
			name:         "whitespace only fallback",
			themeStore:   "   ",
			declEnabled:  false,
			expectedMode: serverconst.StoreModeMutable,
			description:  "when theme.store contains only whitespace, treat as empty",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Setup
			runtime := config.GetServerRuntime()
			runtime.Config.Theme.Store = tc.themeStore
			runtime.Config.DeclarativeResources.Enabled = tc.declEnabled

			// Execute
			mode := getThemeStoreMode()

			// Assert
			suite.Equal(tc.expectedMode, mode, tc.description)
		})
	}
}

// TestIsDeclarativeModeEnabled tests the declarative mode check.
func (suite *ConfigTestSuite) TestIsDeclarativeModeEnabled() {
	testCases := []struct {
		name         string
		themeStore   string
		declEnabled  bool
		expectedMode bool
		description  string
	}{
		{
			name:         "declarative mode enabled",
			themeStore:   "declarative",
			declEnabled:  false,
			expectedMode: true,
			description:  "when theme.store is set to 'declarative'",
		},
		{
			name:         "mutable mode",
			themeStore:   "mutable",
			declEnabled:  false,
			expectedMode: false,
			description:  "when theme.store is set to 'mutable'",
		},
		{
			name:         "composite mode",
			themeStore:   "composite",
			declEnabled:  false,
			expectedMode: false,
			description:  "when theme.store is set to 'composite'",
		},
		{
			name:         "fallback to global declarative enabled",
			themeStore:   "",
			declEnabled:  true,
			expectedMode: true,
			description:  "when theme.store is empty and global declarative is enabled",
		},
		{
			name:         "fallback to global mutable disabled",
			themeStore:   "",
			declEnabled:  false,
			expectedMode: false,
			description:  "when theme.store is empty and global declarative is disabled",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Setup
			runtime := config.GetServerRuntime()
			runtime.Config.Theme.Store = tc.themeStore
			runtime.Config.DeclarativeResources.Enabled = tc.declEnabled

			// Execute
			mode := isDeclarativeModeEnabled()

			// Assert
			suite.Equal(tc.expectedMode, mode, tc.description)
		})
	}
}

// TestConfigTestSuite runs the config test suite.
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
