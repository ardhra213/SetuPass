// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package client

// Initialize initializes and returns the client factory interface.
func Initialize() ClientFactoryInterface {
	return newClientFactory()
}
