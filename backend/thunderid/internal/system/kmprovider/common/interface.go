// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package common defines interfaces for key manager providers.
package common

import (
	"context"
	"crypto/tls"
)

// ConfigCryptoProvider provides symmetric encryption and decryption functionality
// using statically configured keys.
type ConfigCryptoProvider interface {
	Encrypt(ctx context.Context, content []byte) ([]byte, error)
	Decrypt(ctx context.Context, content []byte) ([]byte, error)
}

// TLSConfigProvider provides TLSMaterial for a key reference.
type TLSConfigProvider interface {
	GetTLSMaterial(ctx context.Context) (*TLSMaterial, error)
}

// TLSMaterial holds the TLS certificate material for a key reference.
type TLSMaterial struct {
	Certificate tls.Certificate
	MinVersion  uint16
}
