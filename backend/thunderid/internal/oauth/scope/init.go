// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package scope

// Initialize initializes and returns a new scope validator.
func Initialize() ScopeValidatorInterface {
	return newAPIScopeValidator()
}
