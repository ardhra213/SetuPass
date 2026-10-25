// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package formatter provides formatting for events.
package formatter

import (
	"encoding/json"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// JSONFormatter formats events as JSON.
type JSONFormatter struct{}

var _ FormatterInterface = (*JSONFormatter)(nil)

// NewJSONFormatter creates a new JSON formatter.
func newJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format formats an event as JSON.
func (jf *JSONFormatter) Format(evt *providers.Event) ([]byte, error) {
	return json.Marshal(evt)
}

// GetName returns the name of this formatter.
func (jf *JSONFormatter) GetName() string {
	return "JSONFormatter"
}
