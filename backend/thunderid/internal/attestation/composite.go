// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attestation

import (
	"context"

	"github.com/thunder-id/thunderid/internal/system/log"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// compositeVerifier routes attestation verification to the platform provider that matches the
// application's configuration. An application configures exactly one platform, so the platform is
// determined by which sub-config is present.
type compositeVerifier struct {
	android providers.AttestationProvider
	apple   providers.AttestationProvider
	logger  *log.Logger
}

// newCompositeVerifier creates a platform-dispatching attestation provider.
func newCompositeVerifier(android, apple providers.AttestationProvider) providers.AttestationProvider {
	return &compositeVerifier{
		android: android,
		apple:   apple,
		logger:  log.GetLogger().With(log.String(log.LoggerKeyComponentName, "AttestationVerifier")),
	}
}

// Verify dispatches to the Android (Play Integrity) or Apple (App Attest) verifier based on the
// configured platform. A configuration with no platform set is an operational error.
func (c *compositeVerifier) Verify(ctx context.Context, cfg *providers.AttestationConfig, token string) (
	bool, *tidcommon.ServiceError) {
	switch {
	case cfg != nil && cfg.Android != nil:
		return c.android.Verify(ctx, cfg, token)
	case cfg != nil && cfg.Apple != nil:
		return c.apple.Verify(ctx, cfg, token)
	default:
		c.logger.Error(ctx, "Attestation requested without a platform configuration")
		return false, &tidcommon.InternalServerError
	}
}
