// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package par

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// requestURIRandomBytes is the number of random bytes for the request URI (32 bytes = 256 bits).
const requestURIRandomBytes = 32

// parStoreInterface defines the interface for PAR request storage.
// Implementations operate on opaque random keys; the request_uri URN prefix is
// added and stripped by the service layer.
type parStoreInterface interface {
	Store(ctx context.Context, request pushedAuthorizationRequest, expirySeconds int64) (string, error)
	Consume(ctx context.Context, randomKey string) (pushedAuthorizationRequest, bool, error)
}

// parRequestStore is the runtime-store-backed implementation of parStoreInterface.
type parRequestStore struct {
	storeProvider providers.RuntimeStoreProvider
}

// newPARRequestStore creates a new runtime-store-backed PAR request store.
func newPARRequestStore(storeProvider providers.RuntimeStoreProvider) parStoreInterface {
	return &parRequestStore{
		storeProvider: storeProvider,
	}
}

// Store persists a pushed authorization request and returns the generated random key.
func (s *parRequestStore) Store(
	ctx context.Context, request pushedAuthorizationRequest, expirySeconds int64,
) (string, error) {
	randomKey, err := generateRandomKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate request URI: %w", err)
	}

	data, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal PAR request: %w", err)
	}

	err = s.storeProvider.Put(ctx, providers.NamespacePAR, randomKey, data, expirySeconds)
	if err != nil {
		return "", fmt.Errorf("failed to store PAR request: %w", err)
	}

	return randomKey, nil
}

// Consume atomically retrieves and deletes a pushed authorization request from the store.
// Returns the request, a boolean indicating if found, and any error.
func (s *parRequestStore) Consume(
	ctx context.Context, randomKey string,
) (pushedAuthorizationRequest, bool, error) {
	data, err := s.storeProvider.Take(ctx, providers.NamespacePAR, randomKey)
	if err != nil {
		return pushedAuthorizationRequest{}, false, fmt.Errorf("failed to retrieve PAR request: %w", err)
	}
	if data == nil {
		return pushedAuthorizationRequest{}, false, nil
	}

	var request pushedAuthorizationRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return pushedAuthorizationRequest{}, false, fmt.Errorf("failed to unmarshal PAR request: %w", err)
	}
	return request, true, nil
}

// generateRandomKey generates a cryptographically random key for the request URI.
func generateRandomKey() (string, error) {
	b := make([]byte, requestURIRandomBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
