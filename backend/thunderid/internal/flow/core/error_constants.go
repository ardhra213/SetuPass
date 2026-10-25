// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// Define core flow errors

// ErrExecutorPrerequisiteNotMet is returned when an executor prerequisite is not met.
var ErrExecutorPrerequisiteNotMet = tidcommon.ServiceError{
	Type: tidcommon.ClientErrorType,
	Code: "FLC-1001",
	Error: tidcommon.I18nMessage{
		Key:          "error.flow.core.executor_prerequisite_not_met",
		DefaultValue: "A prerequisite for the executor was not met",
	},
	ErrorDescription: tidcommon.I18nMessage{
		Key: "error.flow.core.executor_prerequisite_not_met_description",
		DefaultValue: "One or more prerequisites required for the executor were not satisfied. " +
			"Please check the inputs and try again.",
	},
}

// ErrInvalidActionProvided is returned when an invalid action is provided in a prompt node.
var ErrInvalidActionProvided = tidcommon.ServiceError{
	Type: tidcommon.ClientErrorType,
	Code: "FLC-1002",
	Error: tidcommon.I18nMessage{
		Key:          "error.flow.core.prompt_invalid_action",
		DefaultValue: "Invalid action provided",
	},
	ErrorDescription: tidcommon.I18nMessage{
		Key:          "error.flow.core.prompt_invalid_action_description",
		DefaultValue: "The action provided is not valid for the current flow step",
	},
}
