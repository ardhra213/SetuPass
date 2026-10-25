// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package dpop

import (
	oauthconfig "github.com/thunder-id/thunderid/internal/oauth/config"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/jti"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize builds a DPoP verifier from the OAuth config, using the provided
// JTI replay-cache store.
func Initialize(cfg oauthconfig.Config, jtiStore jti.JTIStoreInterface,
	runtimeCryptoProvider providers.RuntimeCryptoProvider) VerifierInterface {
	dpopCfg := cfg.OAuth.DPoP

	return newVerifier(
		jtiStore,
		dpopCfg.AllowedAlgs,
		dpopCfg.IatWindow,
		dpopCfg.Leeway,
		dpopCfg.MaxJTILength,
		runtimeCryptoProvider,
	)
}
