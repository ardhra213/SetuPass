// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authz

import (
	"context"
	"net/url"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/authz/requestvalidator"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/resourceindicators"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// AuthorizationValidatorInterface defines the interface for validating OAuth2 authorization requests.
type AuthorizationValidatorInterface interface {
	validateInitialAuthorizationRequest(ctx context.Context, msg *OAuthMessage, oauthApp *providers.OAuthClient) (
		bool, string, string)
}

// authorizationValidator implements the AuthorizationValidatorInterface for validating OAuth2 authorization requests.
type authorizationValidator struct{}

// newAuthorizationValidator creates a new instance of authorizationValidator.
func newAuthorizationValidator() AuthorizationValidatorInterface {
	return &authorizationValidator{}
}

// validateInitialAuthorizationRequest validates the initial authorization request parameters.
func (av *authorizationValidator) validateInitialAuthorizationRequest(ctx context.Context, msg *OAuthMessage,
	oauthApp *providers.OAuthClient) (bool, string, string) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "AuthorizationValidator"))

	queryParams := url.Values(msg.RequestQueryParams)
	clientID := queryParams.Get(constants.RequestParamClientID)
	redirectURI := queryParams.Get(constants.RequestParamRedirectURI)

	if clientID == "" {
		return false, constants.ErrorInvalidRequest, "Missing client_id parameter"
	}

	if err := oauthApp.ValidateRedirectURI(ctx, redirectURI); err != nil {
		logger.Debug(ctx, "Validation failed for redirect URI", log.Error(err))
		return false, constants.ErrorInvalidRequest, "Invalid redirect URI"
	}

	// All subsequent validation errors can be sent to the client application via redirect.
	// The /authorize endpoint does not accept a DPoP header (proofs are bound at /par
	// or /token), so dpopHeaderJkt is always empty here.
	errCode, errMsg := requestvalidator.ValidateAuthorizationRequestParams(msg.RequestQueryParams, oauthApp, "")
	if errCode != "" {
		return true, errCode, errMsg
	}

	if errResp := resourceindicators.ValidateResourceURIs(msg.Resources); errResp != nil {
		return true, errResp.Error, errResp.ErrorDescription
	}

	return false, "", ""
}
