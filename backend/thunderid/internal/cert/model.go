// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cert

// Certificate represents a certificate structure in the system.
type Certificate struct {
	ID      string
	RefType CertificateReferenceType
	RefID   string
	Type    CertificateType
	Value   string
}
