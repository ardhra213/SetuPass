// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package interceptor

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// ErrorInterceptorFailed defines the error for interceptor validation failures.
var ErrorInterceptorFailed = tidcommon.ServiceError{
	Code: "ICS-1001",
	Type: tidcommon.ClientErrorType,
	Error: tidcommon.I18nMessage{
		Key:          "error.interceptor.failed",
		DefaultValue: "Interceptor validation failed",
	},
	ErrorDescription: tidcommon.I18nMessage{
		Key:          "error.interceptor.failed_description",
		DefaultValue: "A flow interceptor rejected the request",
	},
}

// ErrorChallengeTokenInvalid defines the error when a challenge token is not provided.
var ErrorChallengeTokenInvalid = tidcommon.ServiceError{
	Code: "ICS-1002",
	Type: tidcommon.ClientErrorType,
	Error: tidcommon.I18nMessage{
		Key:          "error.interceptor.challenge_token_invalid",
		DefaultValue: "Invalid challenge token",
	},
	ErrorDescription: tidcommon.I18nMessage{
		Key:          "error.interceptor.challenge_token_invalid_description",
		DefaultValue: "The challenge token is missing or invalid",
	},
}

// ErrorCaptchaInvalid defines the error when a captcha token fails verification.
var ErrorCaptchaInvalid = tidcommon.ServiceError{
	Code: "ICS-1003",
	Type: tidcommon.ClientErrorType,
	Error: tidcommon.I18nMessage{
		Key:          "error.interceptor.captcha_invalid",
		DefaultValue: "Invalid captcha",
	},
	ErrorDescription: tidcommon.I18nMessage{
		Key:          "error.interceptor.captcha_invalid_description",
		DefaultValue: "The captcha token could not be verified",
	},
}
