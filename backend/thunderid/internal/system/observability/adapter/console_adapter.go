// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package adapter

import (
	"fmt"
	"os"
	"sync"
)

// ConsoleAdapter writes events to stdout/stderr.
type consoleAdapter struct {
	mu     sync.Mutex
	closed bool
}

var _ OutputAdapterInterface = (*consoleAdapter)(nil)

// newConsoleAdapter creates a new console-based output adapter.
func newConsoleAdapter() *consoleAdapter {
	return &consoleAdapter{
		closed: false,
	}
}

// Write writes data to stdout.
func (ca *consoleAdapter) Write(data []byte) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if ca.closed {
		return fmt.Errorf("console adapter is closed")
	}

	_, err := os.Stdout.Write(append(data, '\n'))
	return err
}

// Flush is a no-op for console adapter as stdout is unbuffered.
func (ca *consoleAdapter) Flush() error {
	return nil
}

// Close closes the console adapter.
func (ca *consoleAdapter) Close() error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.closed = true
	return nil
}

// GetName returns the name of this adapter.
func (ca *consoleAdapter) GetName() string {
	return "ConsoleAdapter"
}
