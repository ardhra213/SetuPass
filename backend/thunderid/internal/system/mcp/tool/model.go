// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package tool provides common models and utilities for MCP tools.
package tool

// IDInput represents a generic input that requires only a resource ID.
type IDInput struct {
	ID string `json:"id" jsonschema:"The unique identifier of the resource"`
}

// PaginationInput represents common input for list operations.
type PaginationInput struct {
	Limit  int `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default: 20)"`
	Offset int `json:"offset,omitempty" jsonschema:"Offset for pagination (default: 0)"`
}
