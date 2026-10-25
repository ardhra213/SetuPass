// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package defaultkm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncryptionService_Encrypt_NoDefaultKey covers lines 48-50: when the
// configCryptoService has no default key, Encrypt should return an error.
func TestEncryptionService_Encrypt_NoDefaultKey(t *testing.T) {
	es := &configCryptoService{
		defaultKeyID: "",
		keys:         map[string][]byte{},
	}
	_, err := es.Encrypt(context.Background(), []byte("plaintext"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "default encryption key not found")
}

// TestEncryptionService_Encrypt_InvalidKeySize covers lines 54-56: when the
// stored key has an invalid AES size, the underlying Encrypt call returns an
// error that the service propagates.
func TestEncryptionService_Encrypt_InvalidKeySize(t *testing.T) {
	es := &configCryptoService{
		defaultKeyID: "bad-key",
		keys:         map[string][]byte{"bad-key": {0x01}}, // 1 byte — invalid AES key length
	}
	_, err := es.Encrypt(context.Background(), []byte("plaintext"))
	require.Error(t, err)
}
