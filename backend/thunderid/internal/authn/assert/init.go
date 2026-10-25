// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package assert

// Initialize initializes the auth assert generator.
func Initialize() AuthAssertGeneratorInterface {
	return newAuthAssertGenerator()
}
