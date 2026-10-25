// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package jwe

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// Client errors for JWE service
var (
	// ErrorDecodingJWE is the error returned when decoding the JWE token fails.
	ErrorDecodingJWE = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "JWE-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.jweservice.decoding_jwe_error",
			DefaultValue: "JWE decode error",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.jweservice.decoding_jwe_error_description",
			DefaultValue: "Error occurred while decoding JWE token",
		},
	}

	// ErrorJWEDecryptionFailed is the error returned when the JWE token decryption fails.
	ErrorJWEDecryptionFailed = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "JWE-1002",
		Error: tidcommon.I18nMessage{
			Key:          "error.jweservice.decryption_failed",
			DefaultValue: "JWE decryption failed",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.jweservice.decryption_failed_description",
			DefaultValue: "Failed to decrypt the JWE token",
		},
	}

	// ErrorUnsupportedJWEAlgorithm is the error returned when the JWE algorithm is unsupported.
	ErrorUnsupportedJWEAlgorithm = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "JWE-1003",
		Error: tidcommon.I18nMessage{
			Key:          "error.jweservice.unsupported_algorithm",
			DefaultValue: "Unsupported JWE algorithm",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.jweservice.unsupported_algorithm_description",
			DefaultValue: "The specified JWE algorithm is not supported",
		},
	}

	// ErrorUnsupportedEncryptionAlgorithm is the error returned when the encryption algorithm is unsupported.
	ErrorUnsupportedEncryptionAlgorithm = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "JWE-1004",
		Error: tidcommon.I18nMessage{
			Key:          "error.jweservice.unsupported_encryption_algorithm",
			DefaultValue: "Unsupported encryption algorithm",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.jweservice.unsupported_encryption_algorithm_description",
			DefaultValue: "The specified encryption algorithm is not supported",
		},
	}
)
