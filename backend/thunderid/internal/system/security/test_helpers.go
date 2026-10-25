// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package security

import (
	"context"
	"testing"
)

// NewSecurityContextForTest creates a new immutable SecurityContext.
// Used for testing purposes.
func NewSecurityContextForTest(userID, ouID, token string,
	permissions []string, attributes map[string]interface{}) *SecurityContext {
	if !testing.Testing() {
		panic("only for tests!")
	}
	return &SecurityContext{
		subject:     userID,
		ouID:        ouID,
		token:       token,
		permissions: permissions,
		attributes:  attributes,
	}
}

// WithSecurityContextTest adds security context to the request context.
// Used for testing purposes.
func WithSecurityContextTest(ctx context.Context, authCtx *SecurityContext) context.Context {
	if !testing.Testing() {
		panic("only for tests!")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, securityContextKey, authCtx)
}
