// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CredentialTestSuite struct {
	suite.Suite
}

func TestCredentialTestSuite(t *testing.T) {
	suite.Run(t, new(CredentialTestSuite))
}

func (s *CredentialTestSuite) TestIsCredential_StringProperty() {
	schema, err := CompileSchema(json.RawMessage(`{
		"password": {"type": "string", "credential": true},
		"email": {"type": "string"}
	}`))
	s.Require().NoError(err)

	s.Require().True(schema.properties["password"].isCredential())
	s.Require().False(schema.properties["email"].isCredential())
}

func (s *CredentialTestSuite) TestIsCredential_NumberProperty() {
	schema, err := CompileSchema(json.RawMessage(`{
		"pin": {"type": "number", "credential": true},
		"age": {"type": "number"}
	}`))
	s.Require().NoError(err)

	s.Require().True(schema.properties["pin"].isCredential())
	s.Require().False(schema.properties["age"].isCredential())
}

func (s *CredentialTestSuite) TestIsCredential_BooleanReturnsFalse() {
	schema, err := CompileSchema(json.RawMessage(`{
		"active": {"type": "boolean"}
	}`))
	s.Require().NoError(err)

	s.Require().False(schema.properties["active"].isCredential())
}

func (s *CredentialTestSuite) TestIsCredential_ObjectReturnsFalse() {
	schema, err := CompileSchema(json.RawMessage(`{
		"address": {"type": "object", "properties": {"city": {"type": "string"}}}
	}`))
	s.Require().NoError(err)

	s.Require().False(schema.properties["address"].isCredential())
}

func (s *CredentialTestSuite) TestIsCredential_ArrayReturnsFalse() {
	schema, err := CompileSchema(json.RawMessage(`{
		"tags": {"type": "array", "items": {"type": "string"}}
	}`))
	s.Require().NoError(err)

	s.Require().False(schema.properties["tags"].isCredential())
}

func (s *CredentialTestSuite) TestGetAttributes_CredentialOnly_ReturnsOnlyCredentialProperties() {
	schema, err := CompileSchema(json.RawMessage(`{
		"password": {"type": "string", "credential": true},
		"apiKey": {"type": "string", "credential": true},
		"email": {"type": "string", "unique": true},
		"age": {"type": "number"},
		"active": {"type": "boolean"}
	}`))
	s.Require().NoError(err)

	attrs := schema.GetAttributes(AttributeFilter{AllowCredential: true})
	names := make([]string, 0, len(attrs))
	for _, a := range attrs {
		names = append(names, a.Attribute)
		s.True(a.Credential, "credential attribute must have Credential=true")
	}
	sort.Strings(names)
	s.Require().Equal([]string{"apiKey", "password"}, names)
}

func (s *CredentialTestSuite) TestGetAttributes_CredentialOnly_EmptyWhenNoCredentials() {
	schema, err := CompileSchema(json.RawMessage(`{
		"email": {"type": "string"},
		"age": {"type": "number"}
	}`))
	s.Require().NoError(err)

	attrs := schema.GetAttributes(AttributeFilter{AllowCredential: true})
	s.Require().Empty(attrs)
}

func (s *CredentialTestSuite) TestCredentialFieldDefaultsFalse() {
	schema, err := CompileSchema(json.RawMessage(`{
		"name": {"type": "string"}
	}`))
	s.Require().NoError(err)

	s.Require().False(schema.properties["name"].isCredential())
}

func (s *CredentialTestSuite) TestCredentialFieldExplicitFalse() {
	schema, err := CompileSchema(json.RawMessage(`{
		"name": {"type": "string", "credential": false}
	}`))
	s.Require().NoError(err)

	s.Require().False(schema.properties["name"].isCredential())
}

func (s *CredentialTestSuite) TestCredentialFieldInvalidValue() {
	_, err := CompileSchema(json.RawMessage(`{
		"password": {"type": "string", "credential": "yes"}
	}`))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "'credential' field must be a boolean")
}

func (s *CredentialTestSuite) TestCredentialFieldNotAllowedOnBoolean() {
	_, err := CompileSchema(json.RawMessage(`{
		"active": {"type": "boolean", "credential": true}
	}`))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid field 'credential'")
}

func (s *CredentialTestSuite) TestCredentialFieldNotAllowedOnObject() {
	_, err := CompileSchema(json.RawMessage(`{
		"address": {"type": "object", "credential": true, "properties": {"city": {"type": "string"}}}
	}`))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid field 'credential'")
}

func (s *CredentialTestSuite) TestCredentialFieldNotAllowedOnArray() {
	_, err := CompileSchema(json.RawMessage(`{
		"tags": {"type": "array", "credential": true, "items": {"type": "string"}}
	}`))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid field 'credential'")
}
