// Copyright 2025-2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"errors"
)

// ApplicationType identifies the platform/client class of an application. It is the canonical
// discriminator the backend uses to apply type-specific behavior (Flow Secret issuance, direct
// flow initiation, and so on). The free-form Template string remains display metadata only.
type ApplicationType string

// Supported application types.
const (
	ApplicationTypeBrowser   ApplicationType = "browser"
	ApplicationTypeFullStack ApplicationType = "fullstack"
	ApplicationTypeMobile    ApplicationType = "mobile"
	ApplicationTypeM2M       ApplicationType = "m2m"
	ApplicationTypeMCP       ApplicationType = "mcp"
	ApplicationTypeCustom    ApplicationType = "custom"
)

// ApplicationNotFoundError is the error returned when an application is not found.
var ApplicationNotFoundError error = errors.New("application not found")

// ApplicationDataCorruptedError is the error returned when application data is corrupted.
var ApplicationDataCorruptedError error = errors.New("application data is corrupted")

// Constants for MCP tool defaults
var (
	// DefaultUserAttributes are the standard user attributes for application templates.
	DefaultUserAttributes = []string{
		"email", "name", "given_name", "family_name",
		"profile", "picture", "phone_number", "address", "created_at",
	}
	// DefaultScopes are the standard OAuth scopes for application templates.
	DefaultScopes = []string{"openid", "profile", "email"}
)
