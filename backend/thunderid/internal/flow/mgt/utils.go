// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package flowmgt

import (
	"regexp"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// handleFormatRegex matches valid handle format:
// - starts with lowercase letter or digit
// - contains only lowercase letters, digits, underscores, or dashes
// - ends with lowercase letter or digit
var handleFormatRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*[a-z0-9]$|^[a-z0-9]$`)

// isValidHandleFormat validates that the handle follows the required format.
func isValidHandleFormat(handle string) bool {
	return handleFormatRegex.MatchString(handle)
}

// isValidFlowType checks if the provided flow type is valid.
func isValidFlowType(flowType providers.FlowType) bool {
	for _, valid := range providers.ValidFlowTypes {
		if flowType == valid {
			return true
		}
	}
	return false
}
