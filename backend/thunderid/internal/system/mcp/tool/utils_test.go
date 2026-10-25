// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package tool

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

type TestStruct struct {
	ID          string   `json:"id" jsonschema:"The ID"`
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags,omitempty"`
	Description string   `json:"description,omitempty"`
}

func (suite *UtilsTestSuite) TestGenerateSchema() {
	schema := GenerateSchema[TestStruct]()

	require.NotNil(suite.T(), schema)
	assert.Equal(suite.T(), "object", schema.Type)
	require.NotNil(suite.T(), schema.Properties)
	assert.Contains(suite.T(), schema.Properties, "id")
	assert.Contains(suite.T(), schema.Properties, "name")
}

func (suite *UtilsTestSuite) TestWithEnum() {
	schema := GenerateSchema[TestStruct](
		WithEnum("", "type", []string{"A", "B"}),
	)

	prop := schema.Properties["type"]
	assert.NotNil(suite.T(), prop)
	assert.Equal(suite.T(), []any{"A", "B"}, prop.Enum)
}

func (suite *UtilsTestSuite) TestWithDefault() {
	schema := GenerateSchema[TestStruct](
		WithDefault("", "name", "default-name"),
	)

	prop := schema.Properties["name"]
	assert.NotNil(suite.T(), prop)
	// Default is stored as json.RawMessage bytes
	expectedJSON, _ := json.Marshal("default-name")
	assert.Equal(suite.T(), expectedJSON, []byte(prop.Default))
}

func (suite *UtilsTestSuite) TestWithRequired() {
	schema := GenerateSchema[TestStruct](
		WithRequired("", "name", "description"),
	)

	assert.Contains(suite.T(), schema.Required, "name")
	assert.Contains(suite.T(), schema.Required, "description")
}

func (suite *UtilsTestSuite) TestWithRemove() {
	schema := GenerateSchema[TestStruct](
		WithRemove("", "id"),
	)

	assert.Nil(suite.T(), schema.Properties["id"])
}
