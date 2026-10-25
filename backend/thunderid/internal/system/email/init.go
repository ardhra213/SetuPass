// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package email

// Initialize creates and returns the configured email client.
func Initialize() (EmailClientInterface, error) {
	return NewSMTPClientFromConfig()
}
