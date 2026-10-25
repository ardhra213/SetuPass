// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package template

import "context"

// templateStoreInterface defines the interface for template store operations.
type templateStoreInterface interface {
	GetTemplate(ctx context.Context, id string) (*TemplateDTO, error)

	GetTemplateByScenario(ctx context.Context, scenario ScenarioType, tmplType TemplateType) (*TemplateDTO, error)

	ListTemplates(ctx context.Context) ([]*TemplateDTO, error)
}
