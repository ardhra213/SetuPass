// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package declarativeresource provides a generic framework for managing declarative resources.
package declarativeresource

// ResourceConfig defines the configuration for a declarative resource type.
type ResourceConfig struct {
	// ResourceType is the name of the resource type (e.g., "IdentityProvider", "Application")
	ResourceType string

	// DirectoryName is the directory name under declarative_resources/ (e.g., "identity_providers")
	DirectoryName string

	// Parser converts YAML bytes to a DTO
	// Returns the parsed DTO or an error
	Parser func([]byte) (interface{}, error)

	// Validator validates the parsed DTO (optional)
	// Returns an error if validation fails
	Validator func(interface{}) error

	// DependencyValidator validates dependencies (optional)
	// This is called after basic validation to check cross-resource dependencies
	DependencyValidator func(interface{}) error

	// IDExtractor extracts the ID from the DTO
	IDExtractor func(interface{}) string
}
