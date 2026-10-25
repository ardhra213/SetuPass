// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConsentModelTestSuite struct {
	suite.Suite
}

func TestConsentModelTestSuite(t *testing.T) {
	suite.Run(t, new(ConsentModelTestSuite))
}

func (s *ConsentModelTestSuite) TestConsentStatus_IsValid() {
	cases := []struct {
		status ConsentStatus
		valid  bool
	}{
		{ConsentStatusActive, true},
		{ConsentStatusExpired, true},
		{"UNKNOWN", false},
		{"", false},
	}
	for _, tc := range cases {
		s.Equal(tc.valid, tc.status.IsValid())
	}
}

func (s *ConsentModelTestSuite) TestNamespace_IsValid() {
	cases := []struct {
		namespace Namespace
		valid     bool
	}{
		{NamespaceAttribute, true},
		{NamespacePermission, true},
		{"other", false},
		{"", false},
	}
	for _, tc := range cases {
		s.Equal(tc.valid, tc.namespace.IsValid())
	}
}

func (s *ConsentModelTestSuite) TestConsentAuthorizationType_IsValid() {
	cases := []struct {
		authType ConsentAuthorizationType
		valid    bool
	}{
		{AuthorizationTypeAuthorization, true},
		{AuthorizationTypeReAuthorization, true},
		{"OTHER", false},
		{"", false},
	}
	for _, tc := range cases {
		s.Equal(tc.valid, tc.authType.IsValid())
	}
}

func (s *ConsentModelTestSuite) TestConsentAuthorizationStatus_IsValid() {
	cases := []struct {
		status ConsentAuthorizationStatus
		valid  bool
	}{
		{AuthorizationStatusCreated, true},
		{AuthorizationStatusApproved, true},
		{AuthorizationStatusRejected, true},
		{"OTHER", false},
		{"", false},
	}
	for _, tc := range cases {
		s.Equal(tc.valid, tc.status.IsValid())
	}
}
