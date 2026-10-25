// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package consent

// Wire-level discriminator values for ConsentPurposePrompt.Type. These are the strings the UI
// reads to choose between attribute and permission rendering.
const (
	consentPromptTypeAttributes  = "attributes"
	consentPromptTypePermissions = "permissions"
)

// consentSessionData holds the consent session state that is signed into a JWT token.
// It captures the purposes and their elements from the resolve step so that the record step
// can verify that the user's decisions match exactly what was prompted.
type consentSessionData struct {
	// Purposes holds the per-purpose element information from the resolve step
	Purposes []consentSessionPurpose `json:"purposes"`
}

// consentSessionPurpose represents a single purpose's elements within the consent session.
type consentSessionPurpose struct {
	// PurposeName is the unique name of the consent purpose
	PurposeName string `json:"purposeName"`
	// Essential holds the names of mandatory elements for this purpose
	Essential []string `json:"essential"`
	// Optional holds the names of optional elements for this purpose
	Optional []string `json:"optional"`
}
