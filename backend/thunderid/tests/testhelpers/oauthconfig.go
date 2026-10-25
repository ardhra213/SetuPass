// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package testhelpers provides shared test fixtures for backend unit tests.
package testhelpers

import (
	oauthconfig "github.com/thunder-id/thunderid/internal/oauth/config"
	engineconfig "github.com/thunder-id/thunderid/pkg/thunderidengine/config"
)

// OAuthConfig returns a minimal OAuth configuration for unit tests.
func OAuthConfig() oauthconfig.Config {
	return oauthconfig.Config{
		DeploymentID:           "test-deployment",
		RuntimeTransientDBType: "sqlite",
		BaseURL:                "https://thunder.io",
		JWT: engineconfig.JWTConfig{
			Issuer:         "https://thunder.io",
			ValidityPeriod: 3600,
			Audience:       "https://thunder.io",
			Leeway:         60,
		},
		OAuth: engineconfig.OAuthConfig{
			RefreshToken: engineconfig.RefreshTokenConfig{
				RenewOnGrant:   false,
				ValidityPeriod: 86400,
			},
			AuthorizationCode: engineconfig.AuthorizationCodeConfig{
				ValidityPeriod: 300,
			},
			PAR: engineconfig.PARConfig{
				ExpiresIn: 600,
			},
			DPoP: engineconfig.DPoPConfig{
				AllowedAlgs: []string{"ES256"},
			},
			CIBA: engineconfig.CIBAConfig{
				IDTokenHintMaxAgeDays: 30,
			},
		},
		GateClient: engineconfig.GateClientConfig{
			Scheme:    "https",
			Hostname:  "localhost",
			Port:      3000,
			LoginPath: "/login",
			ErrorPath: "/error",
		},
	}
}
