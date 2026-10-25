// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package role

// NewRoleExporterForTest creates a new role exporter for testing purposes.
func NewRoleExporterForTest(
	service RoleServiceInterface, assignmentService RoleAssignmentServiceInterface,
) *roleExporter {
	return newRoleExporter(service, assignmentService)
}
