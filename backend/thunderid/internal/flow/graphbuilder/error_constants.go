// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package graphbuilder builds executable flow graphs from flow definitions.
package graphbuilder

import tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

var (
	errorInvalidFlowData = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "FLG-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.flow.graphbuilder.invalid_flow_data",
			DefaultValue: "Invalid flow data",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.flow.graphbuilder.invalid_flow_data_description",
			DefaultValue: "The flow definition contains invalid data",
		},
	}
	errorGraphBuildFailure = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "FLG-1002",
		Error: tidcommon.I18nMessage{
			Key:          "error.flow.graphbuilder.graph_build_failure",
			DefaultValue: "Graph build failure",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.flow.graphbuilder.graph_build_failure_description",
			DefaultValue: "Failed to build the flow graph",
		},
	}
)
