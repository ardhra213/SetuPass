// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package subscriber provides the subscriber interface and implementations for the analytics system.
package subscriber

import (
	"github.com/thunder-id/thunderid/internal/system/observability/event"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// SubscriberInterface is the interface that all event subscribers must implement.
// Subscribers are now responsible for their own activation and configuration.
type SubscriberInterface interface {
	// GetID returns the unique identifier for this subscriber.
	GetID() string

	// GetCategories returns the categories this subscriber is interested in.
	// Return empty slice or slice containing event.CategoryAll to receive all events.
	GetCategories() []event.EventCategory

	// OnEvent is called when a new event is published.
	// Subscribers are responsible for filtering events they don't want to process.
	OnEvent(evt *providers.Event) error

	// Close is called during shutdown to allow cleanup.
	Close() error

	// IsEnabled checks if the subscriber should be activated based on configuration.
	// The config parameter will be *observability.Config.
	// Returns true if the subscriber should be initialized and activated.
	IsEnabled() bool

	// Initialize sets up the subscriber with the provided configuration.
	// This is called after IsEnabled returns true.
	// Returns error if initialization fails.
	Initialize() error
}
