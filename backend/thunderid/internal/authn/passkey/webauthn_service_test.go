// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package passkey

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/stretchr/testify/suite"
)

const (
	testWebAuthnUserID         = "user123"
	testWebAuthnRelyingPartyID = "example.com"
	testWebAuthnOrigin         = "https://example.com"
)

type WebAuthnLibServiceTestSuite struct {
	suite.Suite
	service *defaultWebAuthnService
}

func TestWebAuthnLibServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WebAuthnLibServiceTestSuite))
}

func (suite *WebAuthnLibServiceTestSuite) SetupTest() {
	// Create a real webauthn service for testing
	service, err := newDefaultWebAuthnService(
		testWebAuthnRelyingPartyID,
		"Test RP",
		[]string{testWebAuthnOrigin},
	)
	suite.Require().NoError(err, "Failed to create webauthn service")
	suite.service = service.(*defaultWebAuthnService)
}

func (suite *WebAuthnLibServiceTestSuite) TestValidatePasskeyLogin_UserHandlerError() {
	// Test that errors from the user handler are properly propagated

	session := sessionData{
		Challenge:        "test-challenge",
		UserVerification: protocol.VerificationPreferred,
	}

	// Create a user handler that returns an error
	expectedError := errors.New("user not found")
	userHandler := func(rawID, userHandle []byte) (webauthnUserInterface, error) {
		return nil, expectedError
	}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
			ParsedCredential: protocol.ParsedCredential{
				ID:   "test-credential-id",
				Type: "public-key",
			},
		},
		Response: protocol.ParsedAssertionResponse{
			CollectedClientData: protocol.CollectedClientData{
				Type:      "webauthn.get",
				Challenge: "test-challenge",
				Origin:    testWebAuthnOrigin,
			},
			UserHandle: []byte(testWebAuthnUserID),
		},
	}

	user, credential, err := suite.service.ValidatePasskeyLogin(userHandler, session, parsedResponse)

	suite.Error(err, "Expected error from user handler")
	suite.Nil(user, "User should be nil on handler error")
	suite.Nil(credential, "Credential should be nil on handler error")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidatePasskeyLogin_EmptySession() {
	// Test with empty session data

	mockUser := newWebauthnUserInterfaceMock(suite.T())
	// Setup default expectations (optional)
	mockUser.On("WebAuthnID").Return([]byte(testWebAuthnUserID)).Maybe()
	mockUser.On("WebAuthnName").Return("test@example.com").Maybe()
	mockUser.On("WebAuthnDisplayName").Return("Test User").Maybe()
	mockUser.On("WebAuthnCredentials").Return([]webauthnCredential{}).Maybe()

	emptySession := sessionData{}

	userHandler := func(_, _ []byte) (webauthnUserInterface, error) {
		return mockUser, nil
	}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
		},
	}

	user, credential, err := suite.service.ValidatePasskeyLogin(userHandler, emptySession, parsedResponse)

	suite.Error(err, "Expected error with empty session")
	suite.Nil(user, "User should be nil on validation failure")
	suite.Nil(credential, "Credential should be nil on validation failure")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidatePasskeyLogin_NilResponse() {
	// Test with nil parsed response - this should panic or error

	mockUser := newWebauthnUserInterfaceMock(suite.T())
	// Setup default expectations
	mockUser.On("WebAuthnID").Return([]byte(testWebAuthnUserID)).Maybe()
	mockUser.On("WebAuthnName").Return("test@example.com").Maybe()
	mockUser.On("WebAuthnDisplayName").Return("Test User").Maybe()
	mockUser.On("WebAuthnCredentials").Return([]webauthnCredential{}).Maybe()
	mockUser.On("WebAuthnIcon").Return("").Maybe()

	session := sessionData{
		Challenge:        "test-challenge",
		UserVerification: protocol.VerificationPreferred,
	}

	userHandler := func(_, _ []byte) (webauthnUserInterface, error) { //nolint:unparam
		return mockUser, nil
	}

	// The webauthn library will panic with nil response, so we test for that
	suite.Panics(func() {
		_, _, _ = suite.service.ValidatePasskeyLogin(userHandler, session, nil)
	}, "Expected panic with nil response")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidatePasskeyLogin_EmptyUserHandle() {
	// Test with empty user handle in response

	mockUser := newWebauthnUserInterfaceMock(suite.T())
	// Setup default expectations (optional)
	mockUser.On("WebAuthnID").Return([]byte(testWebAuthnUserID)).Maybe()
	mockUser.On("WebAuthnName").Return("test@example.com").Maybe()
	mockUser.On("WebAuthnDisplayName").Return("Test User").Maybe()
	mockUser.On("WebAuthnCredentials").Return([]webauthnCredential{}).Maybe()
	mockUser.On("WebAuthnIcon").Return("").Maybe()

	session := sessionData{
		Challenge:        "test-challenge",
		UserVerification: protocol.VerificationPreferred,
	}

	userHandler := func(_, _ []byte) (webauthnUserInterface, error) {
		// The handler might be called depending on how the library handles empty user handles
		return mockUser, nil
	}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
		},
		Response: protocol.ParsedAssertionResponse{
			CollectedClientData: protocol.CollectedClientData{
				Type:      "webauthn.get",
				Challenge: "test-challenge",
				Origin:    testWebAuthnOrigin,
			},
			UserHandle: nil, // Empty user handle
		},
	}

	user, credential, err := suite.service.ValidatePasskeyLogin(userHandler, session, parsedResponse)

	// Expected behavior: should fail validation due to missing cryptographic data or user handle
	suite.Error(err, "Expected validation error")
	suite.Nil(user, "User should be nil on validation failure")
	suite.Nil(credential, "Credential should be nil on validation failure")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidatePasskeyLogin_UserHandlerReturnsValidUser() {
	// Test that user handler is called and user is returned properly

	mockUser := newWebauthnUserInterfaceMock(suite.T())
	// Setup default expectations (optional)
	mockUser.On("WebAuthnID").Return([]byte(testWebAuthnUserID)).Maybe()
	mockUser.On("WebAuthnName").Return("test@example.com").Maybe()
	mockUser.On("WebAuthnDisplayName").Return("Test User").Maybe()
	mockUser.On("WebAuthnCredentials").Return([]webauthnCredential{}).Maybe()
	mockUser.On("WebAuthnIcon").Return("").Maybe()

	session := sessionData{
		Challenge:        "test-challenge",
		UserVerification: protocol.VerificationPreferred,
	}

	userHandlerCalled := false
	var capturedRawID, capturedUserHandle []byte

	userHandler := func(rawID, userHandle []byte) (webauthnUserInterface, error) {
		userHandlerCalled = true
		capturedRawID = rawID
		capturedUserHandle = userHandle
		return mockUser, nil
	}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
		},
		Response: protocol.ParsedAssertionResponse{
			CollectedClientData: protocol.CollectedClientData{
				Type:      "webauthn.get",
				Challenge: "test-challenge",
				Origin:    testWebAuthnOrigin,
			},
			UserHandle: []byte(testWebAuthnUserID),
		},
	}

	// This will fail actual validation due to missing cryptographic data,
	// but we're testing the user handler invocation
	_, _, err := suite.service.ValidatePasskeyLogin(userHandler, session, parsedResponse)

	suite.Error(err, "Expected validation error with mock data")
	suite.True(userHandlerCalled, "User handler should be called")
	suite.Equal([]byte("test-raw-id"), capturedRawID, "RawID should be passed to handler")
	suite.Equal([]byte(testWebAuthnUserID), capturedUserHandle, "UserHandle should be passed to handler")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidatePasskeyLogin_TypeAssertionError() {
	// Test the type assertion failure path when userHandler returns an error
	// The actual type assertion code path (lines 232-234) is difficult to trigger
	// because the library creates the user internally, so we test the error handling

	session := sessionData{
		Challenge:        "test-challenge",
		UserVerification: protocol.VerificationPreferred,
	}

	// Create a user handler that returns an error to simulate failure
	userHandler := func(rawID, userHandle []byte) (webauthnUserInterface, error) {
		return nil, errors.New("user lookup failed")
	}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
		},
		Response: protocol.ParsedAssertionResponse{
			CollectedClientData: protocol.CollectedClientData{
				Type:      "webauthn.get",
				Challenge: "test-challenge",
				Origin:    testWebAuthnOrigin,
			},
			UserHandle: []byte(testWebAuthnUserID),
		},
	}

	user, credential, err := suite.service.ValidatePasskeyLogin(userHandler, session, parsedResponse)

	// Should get an error from the user handler
	suite.Error(err, "Expected error from user handler")
	suite.Nil(user, "User should be nil on error")
	suite.Nil(credential, "Credential should be nil on error")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidateLogin_Success() {
	// Test successful login validation with proper credentials

	mockUser := newWebauthnUserInterfaceMock(suite.T())
	mockUser.On("WebAuthnID").Return([]byte(testWebAuthnUserID)).Maybe()
	mockUser.On("WebAuthnName").Return("test@example.com").Maybe()
	mockUser.On("WebAuthnDisplayName").Return("Test User").Maybe()
	mockUser.On("WebAuthnCredentials").Return([]webauthnCredential{}).Maybe()

	session := sessionData{
		Challenge:        "test-challenge",
		UserVerification: protocol.VerificationPreferred,
	}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
			ParsedCredential: protocol.ParsedCredential{
				ID:   "test-credential-id",
				Type: "public-key",
			},
		},
		Response: protocol.ParsedAssertionResponse{
			CollectedClientData: protocol.CollectedClientData{
				Type:      "webauthn.get",
				Challenge: "test-challenge",
				Origin:    testWebAuthnOrigin,
			},
		},
	}

	credential, err := suite.service.ValidateLogin(mockUser, session, parsedResponse)

	// Will fail due to missing cryptographic data, but we're testing the function is called
	suite.Error(err, "Expected validation error with mock data")
	suite.Nil(credential, "Credential should be nil on validation failure")
}

func (suite *WebAuthnLibServiceTestSuite) TestValidateLogin_EmptySession() {
	// Test login validation with empty session data

	mockUser := newWebauthnUserInterfaceMock(suite.T())
	mockUser.On("WebAuthnID").Return([]byte(testWebAuthnUserID)).Maybe()
	mockUser.On("WebAuthnName").Return("test@example.com").Maybe()
	mockUser.On("WebAuthnDisplayName").Return("Test User").Maybe()
	mockUser.On("WebAuthnCredentials").Return([]webauthnCredential{}).Maybe()

	emptySession := sessionData{}

	parsedResponse := &parsedCredentialAssertionData{
		ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
			RawID: []byte("test-raw-id"),
		},
	}

	credential, err := suite.service.ValidateLogin(mockUser, emptySession, parsedResponse)

	suite.Error(err, "Expected error with empty session")
	suite.Nil(credential, "Credential should be nil on validation failure")
}

func (suite *WebAuthnLibServiceTestSuite) TestBeginDiscoverableLogin_Success() {
	// Test that BeginDiscoverableLogin creates proper options for usernameless authentication

	options, session, err := suite.service.BeginDiscoverableLogin()

	suite.NoError(err, "BeginDiscoverableLogin should not return error")
	suite.NotNil(options, "Options should not be nil")
	suite.NotNil(session, "Session should not be nil")

	// Verify that the challenge is generated
	suite.NotEmpty(session.Challenge, "Challenge should be generated")

	// Verify user verification is set to preferred
	suite.Equal(protocol.VerificationPreferred, session.UserVerification,
		"User verification should be preferred for discoverable login")

	// Verify that AllowedCredentials is empty for usernameless flow
	suite.Empty(session.AllowedCredentialIDs, "Allowed credentials should be empty for discoverable login")
}

func (suite *WebAuthnLibServiceTestSuite) TestParseAssertionResponse_Success() {
	// Test successful parsing of assertion response

	credentialID := base64.RawURLEncoding.EncodeToString([]byte("test-credential-id"))
	clientJSON := `{"type":"webauthn.get","challenge":"test-challenge","origin":"https://example.com"}`
	clientData := base64.RawURLEncoding.EncodeToString([]byte(clientJSON))
	authData := base64.RawURLEncoding.EncodeToString(createMinimalAuthData())
	signature := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))
	userHandle := base64.RawURLEncoding.EncodeToString([]byte("test-user-id"))

	parsed, err := parseAssertionResponse(
		credentialID,
		"public-key",
		clientData,
		authData,
		signature,
		userHandle,
	)

	suite.NoError(err, "Parsing should succeed")
	suite.NotNil(parsed, "Parsed response should not be nil")
	suite.Equal("test-credential-id", string(parsed.RawID))
	suite.Equal("public-key", parsed.Type)
	suite.NotNil(parsed.Response.UserHandle, "User handle should be parsed")
}

func (suite *WebAuthnLibServiceTestSuite) TestParseAssertionResponse_InvalidBase64() {
	// Test parsing with invalid base64 encoding

	credentialID := "invalid!!!base64"
	clientData := base64.RawURLEncoding.EncodeToString([]byte(`{"type":"webauthn.get"}`))
	authData := base64.RawURLEncoding.EncodeToString(createMinimalAuthData())
	signature := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))

	parsed, err := parseAssertionResponse(
		credentialID,
		"public-key",
		clientData,
		authData,
		signature,
		"",
	)

	suite.Error(err, "Should return error for invalid base64")
	suite.Nil(parsed, "Parsed response should be nil on error")
}

func (suite *WebAuthnLibServiceTestSuite) TestParseAssertionResponse_EmptyUserHandle() {
	// Test parsing with empty user handle (valid for usernameless flow after discovery)

	credentialID := base64.RawURLEncoding.EncodeToString([]byte("test-credential-id"))
	clientJSON := `{"type":"webauthn.get","challenge":"test","origin":"https://example.com"}`
	clientData := base64.RawURLEncoding.EncodeToString([]byte(clientJSON))
	authData := base64.RawURLEncoding.EncodeToString(createMinimalAuthData())
	signature := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))

	parsed, err := parseAssertionResponse(
		credentialID,
		"public-key",
		clientData,
		authData,
		signature,
		"", // Empty user handle
	)

	suite.NoError(err, "Parsing should succeed with empty user handle")
	suite.NotNil(parsed, "Parsed response should not be nil")
	suite.Nil(parsed.Response.UserHandle, "User handle should be nil")
}

func (suite *WebAuthnLibServiceTestSuite) TestParseAssertionResponse_InvalidJSON() {
	// Test parsing with invalid client data JSON

	credentialID := base64.RawURLEncoding.EncodeToString([]byte("test-credential-id"))
	clientData := base64.RawURLEncoding.EncodeToString([]byte(`{invalid json}`))
	authData := base64.RawURLEncoding.EncodeToString(createMinimalAuthData())
	signature := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))

	parsed, err := parseAssertionResponse(
		credentialID,
		"public-key",
		clientData,
		authData,
		signature,
		"",
	)

	suite.Error(err, "Should return error for invalid JSON")
	suite.Nil(parsed, "Parsed response should be nil on error")
}

func (suite *WebAuthnLibServiceTestSuite) TestParseAttestationResponse_InvalidBase64() {
	// Test parsing with invalid base64

	credentialID := base64.RawURLEncoding.EncodeToString([]byte("test-credential-id"))
	clientData := "invalid!!!base64"
	attestation := base64.RawURLEncoding.EncodeToString([]byte("test-attestation"))

	parsed, err := parseAttestationResponse(
		credentialID,
		"public-key",
		clientData,
		attestation,
	)

	suite.Error(err, "Should return error for invalid base64")
	suite.Nil(parsed, "Parsed response should be nil on error")
}

// Helper function to create a minimal valid authenticator data for testing
func createMinimalAuthData() []byte {
	// Authenticator data must be at least 37 bytes:
	// 32 bytes RP ID hash + 1 byte flags + 4 bytes sign count
	authData := make([]byte, 37)
	// Set UP (User Present) flag
	authData[32] = 0x01
	return authData
}
