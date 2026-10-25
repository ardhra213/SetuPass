// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package template

import "fmt"

// Initialize sets up the template service with a file-based store loaded from declarative resources.
func Initialize() (TemplateServiceInterface, error) {
	fileStore := newTemplateFileBasedStore()

	if err := loadDeclarativeResources(fileStore); err != nil {
		return nil, fmt.Errorf("failed to initialize template service: %w", err)
	}

	service := newTemplateService(fileStore)
	return service, nil
}
