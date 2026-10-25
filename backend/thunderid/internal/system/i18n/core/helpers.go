// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package core provides internationalization support.
package core

// GetDefault returns the default value for a given i18n key.
// Returns the value and true if found, empty string and false otherwise.
func GetDefault(key string) (string, bool) {
	val, ok := defaultMessages[key]
	return val, ok
}

// GetAllDefaults returns a copy of all default messages.
func GetAllDefaults() map[string]string {
	result := make(map[string]string, len(defaultMessages))
	for k, v := range defaultMessages {
		result[k] = v
	}
	return result
}

// GetAllKeys returns all registered i18n keys.
func GetAllKeys() []string {
	keys := make([]string, 0, len(defaultMessages))
	for k := range defaultMessages {
		keys = append(keys, k)
	}
	return keys
}
