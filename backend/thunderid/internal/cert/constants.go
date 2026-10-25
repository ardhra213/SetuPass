// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cert

import (
	"errors"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// CertificateType represents the type of certificates in the system.
type CertificateType = providers.CertificateType

const (
	// CertificateTypeJWKS represents a JSON Web Key Set (JWKS) certificate.
	CertificateTypeJWKS = providers.CertificateTypeJWKS
	// CertificateTypeJWKSURI represents a JWKS URI certificate.
	CertificateTypeJWKSURI = providers.CertificateTypeJWKSURI
)

// CertificateReferenceType represents the type of certificate reference in the system.
type CertificateReferenceType string

const (
	// CertificateReferenceTypeIDP represents a certificate reference for an identity provider.
	CertificateReferenceTypeIDP CertificateReferenceType = "IDP"
	// CertificateReferenceTypeOAuthApp represents a certificate reference for an OAuth client.
	CertificateReferenceTypeOAuthApp CertificateReferenceType = "OAUTH_APP"
)

// ErrCertificateNotFound is the error message when a certificate is not found.
var ErrCertificateNotFound = errors.New("certificate not found")
