// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package sysauthz

// Initialize creates and returns a SystemAuthorizationServiceInterface instance.
// This package exposes no HTTP routes and requires no store — it is a pure service.
func Initialize() (SystemAuthorizationServiceInterface, error) {
	return newSystemAuthorizationService(), nil
}
