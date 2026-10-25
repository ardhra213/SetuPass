// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package flowexec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// flowStoreInterface defines the methods for flow context storage operations.
type flowStoreInterface interface {
	StoreFlowContext(ctx context.Context, dbModel FlowContextDB, expirySeconds int64) error
	GetFlowContext(ctx context.Context, executionID string) (*FlowContextDB, error)
	UpdateFlowContext(ctx context.Context, dbModel FlowContextDB) error
	DeleteFlowContext(ctx context.Context, executionID string) error
}

// flowStore adapts a runtime store provider to flow context storage. Flow contexts are stored
// under the flow namespace, keyed by execution ID, as a serialized FlowContextDB.
type flowStore struct {
	store providers.RuntimeStoreProvider
}

// newFlowStore creates a flow context store backed by the given runtime store provider.
func newFlowStore(store providers.RuntimeStoreProvider) flowStoreInterface {
	return &flowStore{store: store}
}

// StoreFlowContext serializes and stores the flow context with the given TTL in seconds.
func (s *flowStore) StoreFlowContext(ctx context.Context, dbModel FlowContextDB, expirySeconds int64) error {
	data, err := json.Marshal(dbModel)
	if err != nil {
		return fmt.Errorf("failed to marshal flow context: %w", err)
	}
	return s.store.Put(ctx, providers.NamespaceFlow, dbModel.ExecutionID, data, expirySeconds)
}

// GetFlowContext retrieves and deserializes the flow context. Returns nil when not found or expired.
func (s *flowStore) GetFlowContext(ctx context.Context, executionID string) (*FlowContextDB, error) {
	data, err := s.store.Get(ctx, providers.NamespaceFlow, executionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow context: %w", err)
	}
	if data == nil {
		return nil, nil
	}

	var dbModel FlowContextDB
	if err := json.Unmarshal(data, &dbModel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal flow context: %w", err)
	}
	return &dbModel, nil
}

// UpdateFlowContext serializes and updates the stored flow context, preserving its TTL.
func (s *flowStore) UpdateFlowContext(ctx context.Context, dbModel FlowContextDB) error {
	data, err := json.Marshal(dbModel)
	if err != nil {
		return fmt.Errorf("failed to marshal flow context: %w", err)
	}
	return s.store.Update(ctx, providers.NamespaceFlow, dbModel.ExecutionID, data)
}

// DeleteFlowContext removes the flow context.
func (s *flowStore) DeleteFlowContext(ctx context.Context, executionID string) error {
	return s.store.Delete(ctx, providers.NamespaceFlow, executionID)
}
