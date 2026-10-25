// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package consent provides the consent persistence and service layer.
package consent

// Initialize constructs the consent service.
func Initialize(
	inboundClientProvider InboundClientProvider,
) (ConsentServiceInterface, error) {
	return newConsentService(inboundClientProvider)
}
