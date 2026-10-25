// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authz

import (
	"github.com/thunder-id/thunderid/internal/authz/engine"
	"github.com/thunder-id/thunderid/internal/role"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize creates and initializes the authorization service with the RBAC engine.
func Initialize(roleService role.RoleServiceInterface) providers.AuthorizationProvider {
	rbacEngine := engine.NewRBACEngine(roleService)
	return newAuthorizationService(rbacEngine)
}
