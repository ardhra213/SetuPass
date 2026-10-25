// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
}

func TestModelSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

func (suite *ModelTestSuite) TestI18nMessage_String() {
	msg := I18nMessage{
		Key:          "key",
		DefaultValue: "default value",
	}
	assert.Equal(suite.T(), "default value", msg.String())

	msgEmpty := I18nMessage{}
	assert.Equal(suite.T(), "", msgEmpty.String())

	suite.T().Run("substitutes params into the default value", func(t *testing.T) {
		msg := I18nMessage{
			DefaultValue: "Layout is being used by {{param(count)}} application(s)",
			Params:       map[string]string{"count": "3"},
		}
		assert.Equal(t, "Layout is being used by 3 application(s)", msg.String())
	})
}

func (suite *ModelTestSuite) TestSubstituteParams() {
	suite.T().Run("nil params returns template unchanged", func(t *testing.T) {
		assert.Equal(t, "no {{param(x)}} here", substituteParams("no {{param(x)}} here", nil))
	})

	suite.T().Run("single placeholder", func(t *testing.T) {
		assert.Equal(t, "value cannot be empty for property 'scope'",
			substituteParams("value cannot be empty for property '{{param(property)}}'",
				map[string]string{"property": "scope"}))
	})

	suite.T().Run("multiple placeholders", func(t *testing.T) {
		assert.Equal(t, "property 'x' is not supported for IDP type 'OIDC'",
			substituteParams("property '{{param(property)}}' is not supported for IDP type '{{param(idpType)}}'",
				map[string]string{"property": "x", "idpType": "OIDC"}))
	})

	suite.T().Run("whitespace variations", func(t *testing.T) {
		assert.Equal(t, "a-b-c",
			substituteParams("{{param(a)}}-{{ param(b) }}-{{ param( c ) }}",
				map[string]string{"a": "a", "b": "b", "c": "c"}))
	})

	suite.T().Run("missing param keeps placeholder", func(t *testing.T) {
		assert.Equal(t, "known and {{param(unknown)}}",
			substituteParams("{{param(known)}} and {{param(unknown)}}",
				map[string]string{"known": "known"}))
	})
}

func (suite *ModelTestSuite) TestI18nMessage_IsEmpty() {
	suite.T().Run("empty message", func(t *testing.T) {
		assert.True(t, I18nMessage{}.IsEmpty())
	})

	suite.T().Run("non-empty message", func(t *testing.T) {
		assert.False(t, I18nMessage{Key: "key"}.IsEmpty())
	})

	suite.T().Run("value only without key", func(t *testing.T) {
		assert.True(t, I18nMessage{DefaultValue: "val"}.IsEmpty())
	})
}

func (suite *ModelTestSuite) TestI18nMessage_MarshalJSON() {
	suite.T().Run("substitutes params into defaultValue and keeps params", func(t *testing.T) {
		msg := I18nMessage{
			Key:          "error.layoutservice.layout_in_use_description",
			DefaultValue: "Layout is being used by {{param(count)}} application(s)",
			Params:       map[string]string{"count": "3"},
		}

		data, err := json.Marshal(msg)
		assert.NoError(t, err)

		var out map[string]any
		assert.NoError(t, json.Unmarshal(data, &out))
		assert.Equal(t, "error.layoutservice.layout_in_use_description", out["key"])
		assert.Equal(t, "Layout is being used by 3 application(s)", out["defaultValue"])
		assert.Equal(t, map[string]any{"count": "3"}, out["params"])
	})

	suite.T().Run("omits params and leaves defaultValue unchanged when absent", func(t *testing.T) {
		msg := I18nMessage{Key: "key", DefaultValue: "static message"}

		data, err := json.Marshal(msg)
		assert.NoError(t, err)

		assert.JSONEq(t, `{"key":"key","defaultValue":"static message"}`, string(data))
	})
}

func (suite *ModelTestSuite) TestServiceError_WithParams() {
	base := ServiceError{
		Type: ClientErrorType,
		Code: "RES-1018",
		Error: I18nMessage{
			Key:          "error.resourceservice.cannot_modify_declarative_resource_server",
			DefaultValue: "Cannot modify declarative resource server",
		},
		ErrorDescription: I18nMessage{
			Key: "error.resourceservice.cannot_modify_declarative_resource_server_description",
			DefaultValue: "Resource server {{param(id)}} is defined in declarative " +
				"configuration and cannot be modified",
		},
	}

	err := base.WithParams(map[string]string{"id": "server-1"})

	assert.Equal(suite.T(), map[string]string{"id": "server-1"}, err.ErrorDescription.Params)
	assert.Equal(suite.T(),
		"Resource server server-1 is defined in declarative configuration and cannot be modified",
		err.ErrorDescription.String())
	// The base error must remain untouched.
	assert.Nil(suite.T(), base.ErrorDescription.Params)
}

func (suite *ModelTestSuite) TestCustomServiceError() {
	base := ServiceError{
		Type: ClientErrorType,
		Code: "SSE-4030",
		Error: I18nMessage{
			Key:          "error.unauthorized",
			DefaultValue: "Unauthorized",
		},
		ErrorDescription: I18nMessage{
			Key:          "error.unauthorized_description",
			DefaultValue: "The caller is not authorized",
		},
	}

	suite.T().Run("empty custom description preserves base description", func(t *testing.T) {
		err := CustomServiceError(base, I18nMessage{})
		assert.Equal(t, base.Type, err.Type)
		assert.Equal(t, base.Code, err.Code)
		assert.Equal(t, base.Error, err.Error)
		assert.Equal(t, base.ErrorDescription, err.ErrorDescription)
	})

	suite.T().Run("non-empty custom description overrides base description", func(t *testing.T) {
		customDesc := I18nMessage{
			Key:          "error.custom",
			DefaultValue: "Custom description",
		}
		err := CustomServiceError(base, customDesc)
		assert.Equal(t, customDesc, err.ErrorDescription)
		assert.Equal(t, base.Error, err.Error)
	})
}
