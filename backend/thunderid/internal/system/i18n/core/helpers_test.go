// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefault(t *testing.T) {
	// Setup: Mock DefaultMessages
	origDefaults := defaultMessages
	defer func() { defaultMessages = origDefaults }()

	testKey := "test.key"
	testValue := "Test Value"
	defaultMessages = map[string]string{
		testKey: testValue,
	}

	t.Run("Existing Key", func(t *testing.T) {
		val, ok := GetDefault(testKey)
		assert.True(t, ok)
		assert.Equal(t, testValue, val)
	})

	t.Run("Non-Existing Key", func(t *testing.T) {
		val, ok := GetDefault("non.existent.key")
		assert.False(t, ok)
		assert.Empty(t, val)
	})
}

func TestGetAllDefaults(t *testing.T) {
	// Setup: Mock DefaultMessages
	origDefaults := defaultMessages
	defer func() { defaultMessages = origDefaults }()

	defaultMessages = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	defaults := GetAllDefaults()

	assert.Equal(t, len(defaultMessages), len(defaults))

	// Verify content match
	for k, v := range defaultMessages {
		assert.Equal(t, v, defaults[k])
	}

	// Verify it is a copy
	defaults["new_temp_key"] = "temp_value"
	_, ok := defaultMessages["new_temp_key"]
	assert.False(t, ok, "GetAllDefaults should return a copy, not reference")
}

func TestGetAllKeys(t *testing.T) {
	// Setup: Mock DefaultMessages
	origDefaults := defaultMessages
	defer func() { defaultMessages = origDefaults }()

	defaultMessages = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	keys := GetAllKeys()

	assert.Equal(t, len(defaultMessages), len(keys))

	// Verify all keys are present
	for _, k := range keys {
		_, ok := defaultMessages[k]
		assert.True(t, ok, "Key returned by GetAllKeys should exist in DefaultMessages")
	}
}
