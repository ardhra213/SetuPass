// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

// ApplicationListOutput represents the output for list_applications tool.
type ApplicationListOutput struct {
	TotalCount   int                        `json:"total_count" jsonschema:"Total number of applications available."`
	Applications []BasicApplicationResponse `json:"applications" jsonschema:"List of applications."`
}

// ClientIDInput represents input for client ID-based lookups.
type ClientIDInput struct {
	ClientID string `json:"client_id" jsonschema:"OAuth client ID to search for"`
}
