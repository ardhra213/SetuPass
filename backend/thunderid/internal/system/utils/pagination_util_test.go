// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPaginationLinks_MiddlePage(t *testing.T) {
	links := BuildPaginationLinks("/items", 5, 5, 20, "")
	require.Len(t, links, 4)
	assert.Equal(t, "first", links[0].Rel)
	assert.Equal(t, "/items?offset=0&limit=5", links[0].Href)
	assert.Equal(t, "prev", links[1].Rel)
	assert.Equal(t, "/items?offset=0&limit=5", links[1].Href)
	assert.Equal(t, "next", links[2].Rel)
	assert.Equal(t, "/items?offset=10&limit=5", links[2].Href)
	assert.Equal(t, "last", links[3].Rel)
	assert.Equal(t, "/items?offset=15&limit=5", links[3].Href)
}

func TestBuildPaginationLinks_FirstPage(t *testing.T) {
	links := BuildPaginationLinks("/items", 10, 0, 25, "")
	require.Len(t, links, 2)
	assert.Equal(t, "next", links[0].Rel)
	assert.Equal(t, "last", links[1].Rel)
}

func TestBuildPaginationLinks_LastPage(t *testing.T) {
	links := BuildPaginationLinks("/items", 10, 20, 25, "")
	require.Len(t, links, 2)
	assert.Equal(t, "first", links[0].Rel)
	assert.Equal(t, "prev", links[1].Rel)
}

func TestBuildPaginationLinks_SinglePage(t *testing.T) {
	links := BuildPaginationLinks("/items", 10, 0, 5, "")
	require.Len(t, links, 0)
}

func TestBuildPaginationLinks_ZeroLimit(t *testing.T) {
	links := BuildPaginationLinks("/items", 0, 0, 10, "")
	require.Len(t, links, 0)
}

func TestBuildPaginationLinks_NegativeLimit(t *testing.T) {
	links := BuildPaginationLinks("/items", -1, 0, 10, "")
	require.Len(t, links, 0)
}

func TestBuildPaginationLinks_WithExtraQuery(t *testing.T) {
	links := BuildPaginationLinks("/items", 5, 5, 20, "&include=display")
	require.Len(t, links, 4)
	assert.Equal(t, "/items?offset=0&limit=5&include=display", links[0].Href)
	assert.Equal(t, "/items?offset=0&limit=5&include=display", links[1].Href)
	assert.Equal(t, "/items?offset=10&limit=5&include=display", links[2].Href)
	assert.Equal(t, "/items?offset=15&limit=5&include=display", links[3].Href)
}
