// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- ExtractDisplayValue tests ---

func TestExtractDisplayValue_TopLevel(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	assert.Equal(t, "alice@example.com", ExtractDisplayValue(attrs, "email"))
}

func TestExtractDisplayValue_Nested(t *testing.T) {
	attrs := json.RawMessage(`{"profile":{"fullName":"Alice Smith"}}`)
	assert.Equal(t, "Alice Smith", ExtractDisplayValue(attrs, "profile.fullName"))
}

func TestExtractDisplayValue_NonExistentPath(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	assert.Equal(t, "", ExtractDisplayValue(attrs, "missing.field"))
}

func TestExtractDisplayValue_EmptyAttributes(t *testing.T) {
	assert.Equal(t, "", ExtractDisplayValue(json.RawMessage(`{}`), "email"))
}

func TestExtractDisplayValue_NilAttributes(t *testing.T) {
	assert.Equal(t, "", ExtractDisplayValue(nil, "email"))
}

func TestExtractDisplayValue_InvalidJSON(t *testing.T) {
	assert.Equal(t, "", ExtractDisplayValue(json.RawMessage(`invalid`), "email"))
}

func TestExtractDisplayValue_EmptyPath(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	assert.Equal(t, "", ExtractDisplayValue(attrs, ""))
}

func TestExtractDisplayValue_NumericValue(t *testing.T) {
	attrs := json.RawMessage(`{"age":30}`)
	assert.Equal(t, "30", ExtractDisplayValue(attrs, "age"))
}

func TestExtractDisplayValue_BooleanValue(t *testing.T) {
	attrs := json.RawMessage(`{"active":true}`)
	assert.Equal(t, "", ExtractDisplayValue(attrs, "active"))
}

func TestExtractDisplayValue_DeeplyNested(t *testing.T) {
	attrs := json.RawMessage(`{"a":{"b":{"c":"deep"}}}`)
	assert.Equal(t, "deep", ExtractDisplayValue(attrs, "a.b.c"))
}

func TestExtractDisplayValue_NullValue(t *testing.T) {
	attrs := json.RawMessage(`{"email":null}`)
	assert.Equal(t, "", ExtractDisplayValue(attrs, "email"))
}

func TestExtractDisplayValue_PartialPath(t *testing.T) {
	attrs := json.RawMessage(`{"profile":"not-an-object"}`)
	assert.Equal(t, "", ExtractDisplayValue(attrs, "profile.name"))
}

// --- ResolveDisplay tests ---

func TestResolveDisplay_WithDisplayAttr(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	paths := map[string]string{"employee": "email"}
	assert.Equal(t, "alice@example.com", ResolveDisplay("user-1", "employee", attrs, paths))
}

func TestResolveDisplay_FallbackToID(t *testing.T) {
	attrs := json.RawMessage(`{"name":"Alice"}`)
	paths := map[string]string{"employee": "nonexistent"}
	assert.Equal(t, "user-1", ResolveDisplay("user-1", "employee", attrs, paths))
}

func TestResolveDisplay_NilPaths(t *testing.T) {
	assert.Equal(t, "user-1", ResolveDisplay("user-1", "employee", nil, nil))
}

func TestResolveDisplay_EmptyType(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	paths := map[string]string{"employee": "email"}
	assert.Equal(t, "user-1", ResolveDisplay("user-1", "", attrs, paths))
}

func TestResolveDisplay_NestedPath(t *testing.T) {
	attrs := json.RawMessage(`{"profile":{"fullName":"Alice Smith"}}`)
	paths := map[string]string{"employee": "profile.fullName"}
	assert.Equal(t, "Alice Smith", ResolveDisplay("user-1", "employee", attrs, paths))
}

func TestResolveDisplay_UnknownType(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	paths := map[string]string{"employee": "email"}
	assert.Equal(t, "user-1", ResolveDisplay("user-1", "contractor", attrs, paths))
}

func TestResolveDisplay_EmptyPathValue(t *testing.T) {
	attrs := json.RawMessage(`{"email":"alice@example.com"}`)
	paths := map[string]string{"employee": ""}
	assert.Equal(t, "user-1", ResolveDisplay("user-1", "employee", attrs, paths))
}
