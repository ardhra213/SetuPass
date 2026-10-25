// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package passkey

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// sessionStoreInterface defines the interface for WebAuthn session storage.
type sessionStoreInterface interface {
	storeSession(
		ctx context.Context,
		sessionKey string,
		session *sessionData,
		expirySeconds int64,
	) error
	retrieveSession(ctx context.Context, sessionKey string) (*sessionData, error)
	deleteSession(ctx context.Context, sessionKey string) error
}

// sessionStore adapts a runtime store provider to WebAuthn session storage. Sessions are stored
// under the WebAuthn namespace, keyed by session key, as a serialized sessionData.
type sessionStore struct {
	store providers.RuntimeStoreProvider
}

// newSessionStore creates a WebAuthn session store backed by the given runtime store provider.
func newSessionStore(store providers.RuntimeStoreProvider) sessionStoreInterface {
	return &sessionStore{store: store}
}

// storeSession serializes the WebAuthn session data and stores it with the given TTL.
func (s *sessionStore) storeSession(
	ctx context.Context, sessionKey string, session *sessionData, expirySeconds int64) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal passkey session: %w", err)
	}

	return s.store.Put(ctx, providers.NamespaceWebAuthn, sessionKey, data, expirySeconds)
}

// retrieveSession retrieves the WebAuthn session data. Returns (nil, nil) when the session is
// absent or expired.
func (s *sessionStore) retrieveSession(ctx context.Context, sessionKey string) (*sessionData, error) {
	if sessionKey == "" {
		return nil, nil
	}

	data, err := s.store.Get(ctx, providers.NamespaceWebAuthn, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get passkey session: %w", err)
	}
	if data == nil {
		return nil, nil
	}

	var session sessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal passkey session: %w", err)
	}
	return &session, nil
}

// deleteSession removes the WebAuthn session.
func (s *sessionStore) deleteSession(ctx context.Context, sessionKey string) error {
	if sessionKey == "" {
		return nil
	}

	if err := s.store.Delete(ctx, providers.NamespaceWebAuthn, sessionKey); err != nil {
		return fmt.Errorf("failed to delete passkey session: %w", err)
	}
	return nil
}
