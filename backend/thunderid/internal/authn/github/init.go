// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package github

import (
	authnoauth "github.com/thunder-id/thunderid/internal/authn/oauth"
	syshttp "github.com/thunder-id/thunderid/internal/system/http"
)

// Initialize initializes the GitHub OAuth authentication service.
func Initialize(oauthSvc authnoauth.OAuthAuthnServiceInterface) GithubOAuthAuthnServiceInterface {
	httpClient := syshttp.NewHTTPClient()
	return newGithubOAuthAuthnService(oauthSvc, httpClient)
}
