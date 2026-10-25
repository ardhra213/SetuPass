// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package ciba

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AssertionTestSuite struct {
	suite.Suite
}

func TestAssertionTestSuite(t *testing.T) {
	suite.Run(t, new(AssertionTestSuite))
}

// buildCIBAAssertion constructs a minimal unsigned JWT from a payload map.
func buildCIBAAssertion(payload map[string]interface{}) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payloadBytes, _ := json.Marshal(payload)
	payloadEnc := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return header + "." + payloadEnc + ".sig"
}

func (suite *AssertionTestSuite) TestDecodeAttributesFromAssertion_CIBAAuthReqIDWrongType_ReturnsError() {
	assertion := buildCIBAAssertion(map[string]interface{}{
		"sub":                      "user-1",
		"authorization_request_id": 42,
	})

	_, _, err := decodeAttributesFromAssertion(assertion)
	suite.Error(err)
	suite.Contains(err.Error(), "authorization_request_id")
}

func (suite *AssertionTestSuite) TestDecodeAttributesFromAssertion_MissingAuthorizedPermissions_NoError() {
	// authorized_permissions is optional — absence should not cause an error and the
	// field should default to the empty string.
	assertion := buildCIBAAssertion(map[string]interface{}{
		"sub":                      "user-1",
		"authorization_request_id": "auth-req-123",
	})

	claims, _, err := decodeAttributesFromAssertion(assertion)
	suite.NoError(err)
	suite.Equal("user-1", claims.userID)
	suite.Equal("auth-req-123", claims.authReqID)
	suite.Empty(claims.authorizedPermissions)
}
