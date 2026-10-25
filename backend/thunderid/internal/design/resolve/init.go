// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/application"
	layoutmgt "github.com/thunder-id/thunderid/internal/design/layout/mgt"
	thememgt "github.com/thunder-id/thunderid/internal/design/theme/mgt"
	"github.com/thunder-id/thunderid/internal/system/middleware"
)

// Initialize initializes the design resolve service and registers its routes.
func Initialize(
	mux *http.ServeMux,
	themeMgtService thememgt.ThemeMgtServiceInterface,
	layoutMgtService layoutmgt.LayoutMgtServiceInterface,
	applicationService application.ApplicationServiceInterface,
) DesignResolveServiceInterface {
	designResolveService := newDesignResolveService(themeMgtService, layoutMgtService, applicationService)

	if mux != nil {
		designResolveHandler := newDesignResolveHandler(designResolveService)
		registerRoutes(mux, designResolveHandler)
	}
	return designResolveService
}

// registerRoutes registers the routes for design resolve operations.
func registerRoutes(mux *http.ServeMux, resolveHandler *designResolveHandler) {
	opts := middleware.CORSOptions{
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   middleware.DefaultAllowedHeaders,
		AllowCredentials: true,
		MaxAge:           600,
	}
	mux.HandleFunc(middleware.WithCORS("GET /design/resolve", resolveHandler.HandleResolveRequest, opts))
	mux.HandleFunc(middleware.WithCORS("OPTIONS /design/resolve", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, opts))
}
