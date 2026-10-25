// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

// OIDCScope represents an OpenID Connect scope with its metadata
type OIDCScope struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Claims      []string `json:"claims"`
}
