// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package assert

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// Client errors for authentication assertion operations.
var (
	// ErrorNoAuthenticators is the error returned when no authenticators are provided.
	ErrorNoAuthenticators = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "AST-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.assertservice.no_authenticators",
			DefaultValue: "No authenticators",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.assertservice.no_authenticators_description",
			DefaultValue: "Cannot generate assertion without authenticators",
		},
	}
	// ErrorInvalidAuthenticator is the error returned when authenticator name is invalid.
	ErrorInvalidAuthenticator = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "AST-1002",
		Error: tidcommon.I18nMessage{
			Key:          "error.assertservice.invalid_authenticator",
			DefaultValue: "Invalid authenticator",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.assertservice.invalid_authenticator_description",
			DefaultValue: "Authenticator name cannot be empty",
		},
	}
	// ErrorNilAssuranceContext is the error returned when assurance context is nil.
	ErrorNilAssuranceContext = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "AST-1003",
		Error: tidcommon.I18nMessage{
			Key:          "error.assertservice.nil_assurance_context",
			DefaultValue: "Nil assurance context",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.assertservice.nil_assurance_context_description",
			DefaultValue: "Assurance context cannot be nil for verification",
		},
	}
	// ErrorNoAssuranceRequirements is the error returned when no assurance requirements are specified.
	ErrorNoAssuranceRequirements = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "AST-1004",
		Error: tidcommon.I18nMessage{
			Key:          "error.assertservice.no_assurance_requirements",
			DefaultValue: "No assurance requirements",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.assertservice.no_assurance_requirements_description",
			DefaultValue: "At least one assurance level (AAL or IAL) must be specified for verification",
		},
	}
)
