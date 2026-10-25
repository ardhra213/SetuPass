// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package nextjssdk

// integrateNextJSSDKInput represents the input for the integrate_nextjs_sdk tool.
type integrateNextJSSDKInput struct {
	ServerURL string `json:"server_url,omitempty" jsonschema:"URL of the ThunderID server which is the baseurl for SDK"`
}

// integrateNextJSSDKOutput represents the output for the integrate_nextjs_sdk tool.
type integrateNextJSSDKOutput struct {
	Instructions string `json:"instructions"`
	CodeSnippets string `json:"code_snippets"`
}
