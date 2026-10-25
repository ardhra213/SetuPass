// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package graphbuilder builds executable flow graphs from flow definitions.
package graphbuilder

import (
	"context"

	"github.com/thunder-id/thunderid/internal/flow/core"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// GraphBuilderInterface builds and caches executable flow graphs from flow definitions.
type GraphBuilderInterface interface {
	GetGraph(ctx context.Context, flow *providers.CompleteFlowDefinition) (core.GraphInterface, *tidcommon.ServiceError)
	ValidateGraph(ctx context.Context, flow *providers.CompleteFlowDefinition) *tidcommon.ServiceError
	InvalidateCache(ctx context.Context, flowID string)
}
