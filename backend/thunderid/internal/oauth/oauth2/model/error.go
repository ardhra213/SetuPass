// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

// ErrorResponse represents the OAuth2 error response.
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}
