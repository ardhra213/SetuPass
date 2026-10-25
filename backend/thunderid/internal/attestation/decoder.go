// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attestation

import (
	"context"
	"fmt"

	"cloud.google.com/go/auth/credentials"
	"google.golang.org/api/option"
	playintegrity "google.golang.org/api/playintegrity/v1"
)

// integrityTokenDecoder decodes a Play Integrity token into its plaintext payload by calling
// Google's Play Integrity API. It is an internal seam so the API call can be mocked in tests.
type integrityTokenDecoder interface {
	Decode(ctx context.Context, credentialsJSON, packageName, token string) (
		*playintegrity.TokenPayloadExternal, error)
}

// googlePlayIntegrityDecoder decodes tokens by calling the Google Play Integrity API using the
// application's service account credentials.
type googlePlayIntegrityDecoder struct{}

// newGooglePlayIntegrityDecoder creates a token decoder backed by the Google Play Integrity API.
func newGooglePlayIntegrityDecoder() integrityTokenDecoder {
	return &googlePlayIntegrityDecoder{}
}

// Decode calls the Play Integrity decodeIntegrityToken endpoint for the given package.
func (d *googlePlayIntegrityDecoder) Decode(ctx context.Context, credentialsJSON, packageName, token string) (
	*playintegrity.TokenPayloadExternal, error) {
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		CredentialsJSON: []byte(credentialsJSON),
		Scopes:          []string{playintegrity.PlayintegrityScope},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse play integrity credentials: %w", err)
	}

	svc, err := playintegrity.NewService(ctx, option.WithAuthCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create play integrity client: %w", err)
	}

	resp, err := svc.V1.DecodeIntegrityToken(packageName,
		&playintegrity.DecodeIntegrityTokenRequest{IntegrityToken: token}).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("play integrity decode request failed: %w", err)
	}
	return resp.TokenPayloadExternal, nil
}
