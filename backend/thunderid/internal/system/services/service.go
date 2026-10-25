// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package services provides a way to register and manage HTTP routes for the server.
package services

import "net/http"

// The Route struct represents an HTTP route with its method, path, and handler function.
type Route struct {
	Method      string
	Path        string
	HandlerFunc *http.HandlerFunc
}

// The ServiceInterface struct defines the service that will handle the routes.
type ServiceInterface interface {
	RegisterRoutes(mux *http.ServeMux)
}
