// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package userinfo

import "github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

// UserInfoResponse represents the structured response returned by the
// UserInfo service. It supports JSON, JWS, JWE, and NESTED_JWT response types.
// Only one of JSONBody or JWTBody will be populated based on Type.
type UserInfoResponse struct {
	Type     providers.UserInfoResponseType
	JSONBody map[string]interface{}
	JWTBody  string
}
