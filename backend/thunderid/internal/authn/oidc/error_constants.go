// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// Client errors for OIDC authentication.
var (
	// ErrorInvalidIDToken is the error when the ID token is invalid or malformed.
	ErrorInvalidIDToken = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "AUTH-OIDC-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.authoidcservice.invalid_id_token",
			DefaultValue: "Invalid ID token",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.authoidcservice.invalid_id_token_description",
			DefaultValue: "The ID token is invalid or malformed",
		},
	}
	// ErrorInvalidIDTokenSignature is the error when the ID token signature verification fails.
	ErrorInvalidIDTokenSignature = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "AUTH-OIDC-1002",
		Error: tidcommon.I18nMessage{
			Key:          "error.authoidcservice.invalid_id_token_signature",
			DefaultValue: "Invalid ID token signature",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.authoidcservice.invalid_id_token_signature_description",
			DefaultValue: "The ID token signature verification failed",
		},
	}
)
