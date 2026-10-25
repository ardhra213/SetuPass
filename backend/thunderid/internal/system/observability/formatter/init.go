// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package formatter

// Initialize creates and returns a formatter instance based on the specified format type.
//
// Supported format types:
//   - "json": JSON formatter (default)
//   - Future: "csv", "text", etc.
//
// Parameters:
//   - formatType: The type of formatter to create ("json", etc.)
//
// Returns:
//   - FormatterInterface: The initialized formatter instance
//
// Example:
//
//	formatter := formatter.Initialize("json")
//	data, err := formatter.Format(event)
func Initialize(formatType string) FormatterInterface {
	switch formatType {
	case "json":
		return newJSONFormatter()
	default:
		// Default to JSON formatter
		return newJSONFormatter()
	}
}
