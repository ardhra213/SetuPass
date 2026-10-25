// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attestation

import (
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize creates the platform attestation provider, dispatching to Google Play Integrity for
// Android clients and Apple App Attest for iOS clients based on the application's configuration. It
// returns an error if a platform verifier cannot be constructed, so the caller can fail server
// startup rather than run with a non-functional verifier.
func Initialize(cryptoSvc providers.RuntimeCryptoProvider) (providers.AttestationProvider, error) {
	appleRootPEM := config.GetServerRuntime().Config.Attestation.Apple.RootCertificate
	appAttestVsvc, err := newAppAttestVerifier(appleRootPEM)
	if err != nil {
		return nil, err
	}
	return newCompositeVerifier(
		newPlayIntegrityVerifier(newGooglePlayIntegrityDecoder(), cryptoSvc),
		appAttestVsvc,
	), nil
}
