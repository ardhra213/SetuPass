// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package template provides template rendering functionality for templated content.
package template

import (
	"context"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// TemplateServiceInterface defines the interface for template operations.
type TemplateServiceInterface interface {
	// GetTemplateByScenario retrieves a template by its scenario and template type.
	GetTemplateByScenario(
		ctx context.Context,
		scenario ScenarioType,
		tmplType TemplateType,
	) (*TemplateDTO, *tidcommon.ServiceError)

	// Render renders a template with the provided data.
	Render(
		ctx context.Context,
		scenario ScenarioType,
		tmplType TemplateType,
		data TemplateData,
	) (*RenderedTemplate, *tidcommon.ServiceError)
}
