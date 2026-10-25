// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cryptolib

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

// GenerateSecureToken generates a cryptographically random 32-byte token, hex-encoded (64 chars).
func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// HashToken returns the SHA-256 hex digest of the given token.
// The token itself carries 256 bits of entropy, so no salt is needed.
func HashToken(rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(h[:])
}

// ValidateTokenHash checks whether rawToken hashes to storedHash using constant-time comparison.
func ValidateTokenHash(rawToken, storedHash string) bool {
	expected := HashToken(rawToken)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(storedHash)) == 1
}
