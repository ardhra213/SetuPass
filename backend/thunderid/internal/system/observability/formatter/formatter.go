// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package formatter provides formatter interfaces and implementations for events.
package formatter

import (
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// FormatterInterface is the interface for formatting events into different output formats.
type FormatterInterface interface {
	// Format formats an event into bytes.
	Format(evt *providers.Event) ([]byte, error)

	// GetName returns the name of this formatter.
	GetName() string
}
