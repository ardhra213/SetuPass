// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package export

import declarativeresource "github.com/thunder-id/thunderid/internal/system/declarative_resource"

// ResourceExporterRegistry holds all registered resource exporters.
type resourceExporterRegistry struct {
	exporters map[string]declarativeresource.ResourceExporter
}

// NewResourceExporterRegistry creates a new registry for resource exporters.
func newResourceExporterRegistry() *resourceExporterRegistry {
	return &resourceExporterRegistry{
		exporters: make(map[string]declarativeresource.ResourceExporter),
	}
}

// Register adds a new resource exporter to the registry.
func (r *resourceExporterRegistry) Register(exporter declarativeresource.ResourceExporter) {
	r.exporters[exporter.GetResourceType()] = exporter
}

// Get retrieves a resource exporter by type.
func (r *resourceExporterRegistry) Get(resourceType string) (declarativeresource.ResourceExporter, bool) {
	exporter, exists := r.exporters[resourceType]
	return exporter, exists
}

// GetAll returns all registered exporters.
func (r *resourceExporterRegistry) GetAll() map[string]declarativeresource.ResourceExporter {
	return r.exporters
}
