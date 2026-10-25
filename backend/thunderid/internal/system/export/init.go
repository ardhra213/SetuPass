// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package export provides functionality for exporting various resource configurations.
package export

import (
	"net/http"

	declarativeresource "github.com/thunder-id/thunderid/internal/system/declarative_resource"
	"github.com/thunder-id/thunderid/internal/system/middleware"
)

// Initialize initializes the export service and registers its routes.
func Initialize(mux *http.ServeMux, exporters []declarativeresource.ResourceExporter) ExportServiceInterface {
	// Create parameterizer instance (no longer needs centralized rules)
	parameterizerInstance := newParameterizer(templatingRules{})

	// Create the export service with exporters
	exportService := newExportService(exporters, parameterizerInstance)

	// Create the handler
	exportHandler := newExportHandler(exportService)

	// Register routes
	registerRoutes(mux, exportHandler)

	return exportService
}

func registerRoutes(mux *http.ServeMux, exportHandler *exportHandler) {
	opts := middleware.CORSOptions{
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   middleware.DefaultAllowedHeaders,
		AllowCredentials: true,
		MaxAge:           600,
	}

	// JSON export endpoint
	mux.HandleFunc(middleware.WithCORS("POST /export",
		exportHandler.HandleExportRequest, opts))

	mux.HandleFunc(middleware.WithCORS("OPTIONS /export",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}, opts))
}
