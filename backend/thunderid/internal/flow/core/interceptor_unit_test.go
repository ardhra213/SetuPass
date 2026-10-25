// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"testing"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

	"github.com/stretchr/testify/suite"
)

type InterceptorUnitTestSuite struct {
	suite.Suite
}

func TestInterceptorUnitTestSuite(t *testing.T) {
	suite.Run(t, new(InterceptorUnitTestSuite))
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_GettersReturnSetValues() {
	decl := &interceptorUnit{
		name:    "CaptchaInterceptor",
		mode:    providers.InterceptorModePreNode,
		scope:   providers.InterceptorScopeSelected,
		applyTo: []string{"login-node", "otp-node"},
		properties: map[string]interface{}{
			"threshold": 0.5,
		},
	}

	s.Equal("CaptchaInterceptor", decl.GetName())
	s.Equal(providers.InterceptorModePreNode, decl.GetMode())
	s.Equal(providers.InterceptorScopeSelected, decl.GetScope())
	s.Equal([]string{"login-node", "otp-node"}, decl.GetApplyTo())
	s.Equal(0.5, decl.GetProperties()["threshold"])
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_DefaultZeroValues() {
	decl := &interceptorUnit{}

	s.Empty(decl.GetName())
	s.Empty(decl.GetMode())
	s.Empty(decl.GetScope())
	s.Nil(decl.GetApplyTo())
	s.Nil(decl.GetProperties())
	s.Nil(decl.GetInterceptor())
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_SetName() {
	decl := &interceptorUnit{}
	decl.SetName("RateLimitInterceptor")

	s.Equal("RateLimitInterceptor", decl.GetName())
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_SetMode() {
	decl := &interceptorUnit{}
	decl.SetMode(providers.InterceptorModePostRequest)

	s.Equal(providers.InterceptorModePostRequest, decl.GetMode())
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_SetScope() {
	decl := &interceptorUnit{}
	decl.SetScope(providers.InterceptorScopeAll)

	s.Equal(providers.InterceptorScopeAll, decl.GetScope())
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_SetApplyTo() {
	decl := &interceptorUnit{}
	decl.SetApplyTo([]string{"node-a", "node-b"})

	s.Equal([]string{"node-a", "node-b"}, decl.GetApplyTo())
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_SetProperties() {
	decl := &interceptorUnit{}
	props := map[string]interface{}{"maxAttempts": 5, "window": "60s"}
	decl.SetProperties(props)

	s.Equal(props, decl.GetProperties())
}

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_SetInterceptor() {
	decl := &interceptorUnit{}
	ic := newInterceptor(testInterceptorName, true, 10)
	decl.SetInterceptor(ic)

	s.NotNil(decl.GetInterceptor())
	s.Equal(testInterceptorName, decl.GetInterceptor().GetName())
}

// Interface compliance

func (s *InterceptorUnitTestSuite) TestInterceptorDeclaration_ImplementsInterface() {
	var _ InterceptorUnitInterface = (*interceptorUnit)(nil)
}
