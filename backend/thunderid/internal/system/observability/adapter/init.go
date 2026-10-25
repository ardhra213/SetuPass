// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package adapter

// InitializeConsoleAdapter creates and returns a console adapter that writes formatted events to stdout.
//
// Returns:
//   - OutputAdapterInterface: The initialized console adapter instance
//
// Example:
//
//	adapter := adapter.InitializeConsoleAdapter()
//	err := adapter.Write(formattedData)
func InitializeConsoleAdapter() OutputAdapterInterface {
	return newConsoleAdapter()
}

// InitializeFileAdapter creates and returns a file adapter that writes formatted events to a file
// with optional rotation support.
//
// Parameters:
//   - filePath: The path to the file where events will be written
//
// Returns:
//   - OutputAdapterInterface: The initialized file adapter instance
//   - error: Error if the adapter cannot be created (e.g., invalid path, permission issues)
//
// Example:
//
//	adapter, err := adapter.InitializeFileAdapter("/var/log/observability/events.log")
//	if err != nil {
//	    return err
//	}
//	err = adapter.Write(formattedData)
func InitializeFileAdapter(filePath string) (OutputAdapterInterface, error) {
	return NewFileAdapter(filePath)
}
