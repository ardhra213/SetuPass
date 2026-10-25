// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// InterceptorInterface defines the contract for flow interceptors.
type InterceptorInterface interface {
	// Execute runs the interceptor logic and returns a result.
	Execute(ctx *InterceptorContext) (*common.InterceptorResponse, error)

	// GetName returns the unique name of the interceptor.
	GetName() string

	// IsDefault returns true if this is a default (always enforced) interceptor.
	IsDefault() bool

	// GetPriority returns the execution order within the same mode (lower runs first).
	GetPriority() int

	// GetInputs returns the inputs declared by the interceptor.
	GetInputs() []providers.Input
}

// interceptor represents the basic implementation of an interceptor.
type interceptor struct {
	Name      string
	isDefault bool
	Priority  int
	Inputs    []providers.Input
}

var _ InterceptorInterface = (*interceptor)(nil)

// newInterceptor creates a new instance of interceptor with the given properties.
func newInterceptor(name string, isDefault bool, priority int) InterceptorInterface {
	return &interceptor{
		Name:      name,
		isDefault: isDefault,
		Priority:  priority,
	}
}

// GetName returns the name of the interceptor.
func (i *interceptor) GetName() string {
	return i.Name
}

// IsDefault returns true if this is a default (always enforced) interceptor.
func (i *interceptor) IsDefault() bool {
	return i.isDefault
}

// GetPriority returns the execution order within the same mode (lower runs first).
func (i *interceptor) GetPriority() int {
	return i.Priority
}

// GetInputs returns the inputs declared by the interceptor.
func (i *interceptor) GetInputs() []providers.Input {
	return i.Inputs
}

// Execute runs the interceptor logic and returns a result.
func (i *interceptor) Execute(ctx *InterceptorContext) (*common.InterceptorResponse, error) {
	// Placeholder implementation. Concrete interceptors should override this method.
	return nil, nil
}
