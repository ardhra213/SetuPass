// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package openid4vci

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OID4VCIErrorTestSuite struct {
	suite.Suite
}

func TestOID4VCIErrorTestSuite(t *testing.T) {
	suite.Run(t, new(OID4VCIErrorTestSuite))
}

func (s *OID4VCIErrorTestSuite) TestToOID4VCIError() {
	cases := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"invalid token", ErrInvalidToken, http.StatusUnauthorized, errCodeInvalidToken},
		{"user not found", ErrUserNotFound, http.StatusUnauthorized, errCodeInvalidToken},
		{"invalid dpop", ErrInvalidDPoP, http.StatusUnauthorized, errCodeInvalidDPoPProof},
		{"invalid nonce", ErrInvalidNonce, http.StatusBadRequest, errCodeInvalidNonce},
		{"invalid proof", ErrInvalidProof, http.StatusBadRequest, errCodeInvalidProof},
		{"unsupported", ErrUnsupportedCredential, http.StatusBadRequest, errCodeUnsupportedCredentialType},
		{"invalid request", ErrInvalidRequest, http.StatusBadRequest, errCodeInvalidCredentialRequest},
		{"wrapped proof", fmt.Errorf("decode: %w", ErrInvalidProof), http.StatusBadRequest, errCodeInvalidProof},
		{"unknown", errors.New("boom"), http.StatusInternalServerError, errCodeServerError},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			got := toOID4VCIError(c.err)
			s.Equal(c.status, got.Status)
			s.Equal(c.code, got.Code)
		})
	}
}
