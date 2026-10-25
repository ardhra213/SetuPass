// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package model defines the data structures for the health check module.
package model

// ServerStatus represents an liveliness details of the server.
type ServerStatus struct {
	Status        Status          `json:"status,omitempty"`
	ServiceStatus []ServiceStatus `json:"serviceStatus,omitempty"`
}

// ServiceStatus represents the status of a service in the system.
type ServiceStatus struct {
	ServiceName string `json:"serviceName,omitempty"`
	Status      Status `json:"status,omitempty"`
}

// Status defines the status for service or server.
type Status string

// Status constants represent the possible statuses of a service.
const (
	// StatusUp indicates that the service is operational.
	StatusUp Status = "UP"
	// StatusDown indicates that the service is not operational.
	StatusDown Status = "DOWN"
	// StatusUnknown indicates that the service status is unknown.
	StatusUnknown Status = "UNKNOWN"
)
