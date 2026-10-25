// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package actorprovider

// ApplicationMetadata is the application view composed from inbound-client and actor records.
type ApplicationMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logoUrl,omitempty"`
	URL         string `json:"url,omitempty"`
	TosURI      string `json:"tosUri,omitempty"`
	PolicyURI   string `json:"policyUri,omitempty"`
}
