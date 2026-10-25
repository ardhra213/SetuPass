// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package observability

import (
	"context"

	"github.com/thunder-id/thunderid/internal/system/observability/publisher"
	"github.com/thunder-id/thunderid/internal/system/observability/subscriber"
	engineconfig "github.com/thunder-id/thunderid/pkg/thunderidengine/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// ObservabilityServiceInterface defines the contract for the observability service.
// This interface enables dependency injection and facilitates testing.
// Services that need observability should accept this interface as a parameter.
type ObservabilityServiceInterface interface {
	// PublishEvent publishes an event to the observability system.
	// This is a no-op if observability is disabled.
	// The context carries the request trace ID used for correlated logging.
	PublishEvent(ctx context.Context, evt *providers.Event)

	// IsEnabled returns true if observability is enabled and operational.
	IsEnabled() bool

	// GetConfig returns the current observability configuration.
	GetConfig() *engineconfig.ObservabilityConfig

	// GetPublisher returns the underlying publisher for advanced use cases.
	// Most users should use PublishEvent() instead.
	// Returns nil if observability is disabled.
	GetPublisher() publisher.CategoryPublisherInterface

	// GetActiveSubscribers returns the list of active subscribers.
	// This is useful for testing or querying subscriber state.
	// Returns empty slice if no subscribers are active or observability is disabled.
	GetActiveSubscribers() []subscriber.SubscriberInterface

	// Shutdown gracefully shuts down the observability service.
	// The publisher handles unsubscribing and closing all subscribers.
	Shutdown()
}
