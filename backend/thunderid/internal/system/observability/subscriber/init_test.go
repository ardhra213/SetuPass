// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package subscriber

import (
	"testing"
)

// TestInitialize_RegistryHasSubscribers tests that the registry contains expected subscribers
func TestInitialize_RegistryHasSubscribers(t *testing.T) {
	// Get registered names
	names := GetRegisteredNames()

	// Should have at least the three built-in subscribers
	if len(names) < 3 {
		t.Errorf("Expected at least 3 registered subscribers, got %d", len(names))
	}

	// Check for expected subscriber types
	expectedTypes := map[string]bool{
		"console": false,
		"file":    false,
		"otel":    false,
	}

	for _, name := range names {
		if _, exists := expectedTypes[name]; exists {
			expectedTypes[name] = true
		}
	}

	// Verify all expected types were found
	for typeName, found := range expectedTypes {
		if !found {
			t.Errorf("Expected to find '%s' subscriber in registry", typeName)
		}
	}
}

// TestInitialize_FactoriesCreateInstances tests that all registered factories can create instances
func TestInitialize_FactoriesCreateInstances(t *testing.T) {
	factories := getAllFactories()

	if len(factories) < 3 {
		t.Errorf("Expected at least 3 factories, got %d", len(factories))
	}

	for name, factory := range factories {
		instance := factory()
		if instance == nil {
			t.Errorf("Factory '%s' returned nil instance", name)
			continue
		}

		// Verify instance has required interface methods (they may return zero values until Initialize is called)
		_ = instance.GetID()         // Should not panic
		_ = instance.GetCategories() // Should not panic
		// Note: Can't call IsEnabled() without server runtime being initialized

		// Verify the instance has the SubscriberInterface
		var _ SubscriberInterface = instance
	}
}
