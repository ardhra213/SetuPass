// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package thememgt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ThemeModelTestSuite struct {
	suite.Suite
}

func TestThemeModelTestSuite(t *testing.T) {
	suite.Run(t, new(ThemeModelTestSuite))
}

// Test extractThemeColorInfo - valid light scheme
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_ValidLight() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "light",
		"colorSchemes": {
			"light": {"palette": {"primary": {"main": "#ff7300"}}}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "light", scheme)
	assert.Equal(suite.T(), "#ff7300", color)
}

// Test extractThemeColorInfo - valid dark scheme
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_ValidDark() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "dark",
		"colorSchemes": {
			"dark": {"palette": {"primary": {"main": "#bb86fc"}}}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "dark", scheme)
	assert.Equal(suite.T(), "#bb86fc", color)
}

// Test extractThemeColorInfo - picks the defaultColorScheme when multiple schemes exist
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_MultipleSchemes() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "dark",
		"colorSchemes": {
			"light": {"palette": {"primary": {"main": "#ff7300"}}},
			"dark":  {"palette": {"primary": {"main": "#bb86fc"}}}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "dark", scheme)
	assert.Equal(suite.T(), "#bb86fc", color)
}

// Test extractThemeColorInfo - nil input
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_NilInput() {
	scheme, color := extractThemeColorInfo(nil)

	assert.Equal(suite.T(), "", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - empty input
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_EmptyInput() {
	scheme, color := extractThemeColorInfo(json.RawMessage{})

	assert.Equal(suite.T(), "", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - invalid JSON
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_InvalidJSON() {
	scheme, color := extractThemeColorInfo(json.RawMessage(`{invalid json}`))

	assert.Equal(suite.T(), "", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - missing defaultColorScheme field
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_MissingDefaultColorScheme() {
	themeJSON := json.RawMessage(`{
		"colorSchemes": {
			"light": {"palette": {"primary": {"main": "#ff7300"}}}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - defaultColorScheme references a scheme that doesn't exist
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_SchemeNotInColorSchemes() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "ocean",
		"colorSchemes": {
			"light": {"palette": {"primary": {"main": "#ff7300"}}}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "ocean", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - primary.main is empty string
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_EmptyPrimaryMain() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "light",
		"colorSchemes": {
			"light": {"palette": {"primary": {"main": ""}}}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "light", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - colorSchemes is empty object
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_EmptyColorSchemes() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "light",
		"colorSchemes": {}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "light", scheme)
	assert.Equal(suite.T(), "", color)
}

// Test extractThemeColorInfo - extra fields in JSON are ignored
func (suite *ThemeModelTestSuite) TestExtractThemeColorInfo_ExtraFieldsIgnored() {
	themeJSON := json.RawMessage(`{
		"defaultColorScheme": "light",
		"direction": "ltr",
		"shape": {"borderRadius": 4},
		"colorSchemes": {
			"light": {
				"palette": {
					"primary": {"main": "#ff7300", "light": "#ffa040", "dark": "#b35000"},
					"secondary": {"main": "#9c27b0"},
					"background": {"default": "#f5f5f5", "paper": "#ffffff"}
				}
			}
		}
	}`)

	scheme, color := extractThemeColorInfo(themeJSON)

	assert.Equal(suite.T(), "light", scheme)
	assert.Equal(suite.T(), "#ff7300", color)
}
