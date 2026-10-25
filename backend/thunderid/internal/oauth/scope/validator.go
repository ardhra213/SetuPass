// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package scope provides functionality for validating scopes.
package scope

import "context"

// ScopeError represents an error during scope validation.
type ScopeError struct {
	Error            string
	ErrorDescription string
}

// ScopeValidatorInterface defines the interface for scope validation.
type ScopeValidatorInterface interface {
	ValidateScopes(ctx context.Context, requestedScopes, clientID string) (string, *ScopeError)
}

// apiScopeValidator is the implementation of API scope validation.
type apiScopeValidator struct{}

// newAPIScopeValidator creates a new instance of the apiScopeValidator.
func newAPIScopeValidator() *apiScopeValidator {
	return &apiScopeValidator{}
}

// ValidateScopes validates and filters the requested scopes against the authorized scopes for the application.
func (sv *apiScopeValidator) ValidateScopes(
	ctx context.Context, requestedScopes, clientID string,
) (string, *ScopeError) {
	if requestedScopes == "" {
		return "", nil
	}

	// Return all requested scopes for now.
	return requestedScopes, nil
}
