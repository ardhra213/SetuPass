// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"errors"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// Internal error definitions for template operations.
var (
	// errTemplateNotFound indicates the requested template was not found.
	errTemplateNotFound = errors.New("template not found")
)

// Client errors for template operations.
var (
	// ErrorTemplateNotFound is returned when the requested template does not exist.
	ErrorTemplateNotFound = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "TMP-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.templateservice.template_not_found",
			DefaultValue: "Template not found",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.templateservice.template_not_found_description",
			DefaultValue: "The requested template does not exist for the given scenario",
		},
	}
)
