// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cryptolib

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

func (suite *TokenTestSuite) TestGenerateSecureToken() {
	tok1, err := GenerateSecureToken()
	suite.NoError(err)
	suite.Len(tok1, 64, "token should be 64 hex chars (32 bytes)")

	tok2, err := GenerateSecureToken()
	suite.NoError(err)
	suite.NotEqual(tok1, tok2, "two tokens should not be equal")
}

func (suite *TokenTestSuite) TestHashToken() {
	testCases := []struct {
		name  string
		input string
	}{
		{"ShortInput", "abc123"},
		{"LongInput", "this-is-a-longer-token-string-for-testing"},
		{"EmptyInput", ""},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			hash1 := HashToken(tc.input)
			hash2 := HashToken(tc.input)
			suite.Equal(hash1, hash2, "same input should produce same hash")
			suite.Len(hash1, 64, "SHA-256 hex digest should be 64 chars")
		})
	}

	suite.NotEqual(HashToken("abc"), HashToken("def"), "different inputs should produce different hashes")
}

func (suite *TokenTestSuite) TestValidateTokenHash() {
	tok, err := GenerateSecureToken()
	suite.Require().NoError(err)

	h := HashToken(tok)

	testCases := []struct {
		name       string
		rawToken   string
		storedHash string
		expectOk   bool
	}{
		{"ValidToken", tok, h, true},
		{"WrongToken", "wrong", h, false},
		{"WrongHash", tok, "not-a-valid-hash", false},
		{"EmptyToken", "", h, false},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			suite.Equal(tc.expectOk, ValidateTokenHash(tc.rawToken, tc.storedHash))
		})
	}
}
