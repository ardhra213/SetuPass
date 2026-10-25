// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newServer creates a new MCP server.
func newServer() *mcp.Server {
	// Create the MCP server instance.
	return mcp.NewServer(&mcp.Implementation{
		Name:    "thunderid-mcp",
		Version: "1.0.0",
	}, nil)
}
