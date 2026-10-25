// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attributecache

import (
	"context"
	"fmt"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// attributeCacheStoreInterface defines the interface for the attribute cache store.
type attributeCacheStoreInterface interface {
	// CreateAttributeCache creates a new attribute cache entry in the store.
	CreateAttributeCache(ctx context.Context, id string, data []byte, ttlSeconds int64) error

	// GetAttributeCache retrieves an attribute cache entry by ID from the store.
	GetAttributeCache(ctx context.Context, id string) ([]byte, error)

	// ExtendAttributeCacheTTL extends the TTL of an attribute cache entry in the store.
	ExtendAttributeCacheTTL(ctx context.Context, id string, ttlSeconds int) error

	// DeleteAttributeCache deletes an attribute cache entry by ID from the store.
	DeleteAttributeCache(ctx context.Context, id string) error
}

// attributeCacheStore is the SQL implementation of attributeCacheStoreInterface.
type attributeCacheStore struct {
	store providers.RuntimeStoreProvider
}

// newAttributeCacheStore creates a new instance of attributeCacheStore.
func newAttributeCacheStore(store providers.RuntimeStoreProvider) attributeCacheStoreInterface {
	return &attributeCacheStore{
		store: store,
	}
}

// CreateAttributeCache creates a new attribute cache entry in the database.
func (s *attributeCacheStore) CreateAttributeCache(
	ctx context.Context, id string, data []byte, ttlSeconds int64) error {
	return s.store.Put(ctx, providers.NamespaceAttributeCache, id, data, ttlSeconds)
}

// GetAttributeCache retrieves an attribute cache entry by ID from the database.
func (s *attributeCacheStore) GetAttributeCache(ctx context.Context, id string) ([]byte, error) {
	data, err := s.store.Get(ctx, providers.NamespaceAttributeCache, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute cache: %w", err)
	}
	if data == nil {
		return nil, errAttributeCacheNotFound
	}
	return data, nil
}

// ExtendAttributeCacheTTL extends the TTL of an attribute cache entry in the database.
func (s *attributeCacheStore) ExtendAttributeCacheTTL(ctx context.Context, id string, ttlSeconds int) error {
	return s.store.ExtendTTL(ctx, providers.NamespaceAttributeCache, id, int64(ttlSeconds))
}

// DeleteAttributeCache deletes an attribute cache entry by ID from the database.
func (s *attributeCacheStore) DeleteAttributeCache(ctx context.Context, id string) error {
	return s.store.Delete(ctx, providers.NamespaceAttributeCache, id)
}
