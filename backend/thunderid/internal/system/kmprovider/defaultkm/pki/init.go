// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package pki

// Initialize initializes the PKI service, loading all key/certificate pairs from configuration.
func Initialize() (PKIServiceInterface, error) {
	return newPKIService()
}
