// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package adapter provides output adapter interfaces and implementations.
package adapter

// OutputAdapterInterface is the interface for writing formatted events to various destinations.
type OutputAdapterInterface interface {
	// Write writes formatted event data to the output destination.
	Write(data []byte) error

	// Flush ensures all buffered data is written.
	Flush() error

	// Close closes the adapter and releases resources.
	Close() error

	// GetName returns the name of this adapter.
	GetName() string
}
