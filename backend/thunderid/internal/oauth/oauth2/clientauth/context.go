// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package clientauth

import (
	"context"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

type contextKey string

// OAuthClientKey is the context key for storing authenticated OAuth client info.
var OAuthClientKey contextKey = "oauth_client"

// OAuthClientInfo contains authenticated client information.
type OAuthClientInfo struct {
	ClientID     string
	ClientSecret string
	OAuthApp     *providers.OAuthClient
}

// withOAuthClient adds OAuth client information to the context.
func withOAuthClient(ctx context.Context, client *OAuthClientInfo) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, OAuthClientKey, client)
}

// GetOAuthClient retrieves OAuth client information from the context.
func GetOAuthClient(ctx context.Context) *OAuthClientInfo {
	if ctx == nil {
		return nil
	}

	if client, ok := ctx.Value(OAuthClientKey).(*OAuthClientInfo); ok {
		return client
	}

	return nil
}
