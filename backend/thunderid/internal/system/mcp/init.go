// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package mcp provides MCP (Model Context Protocol) server functionality.
package mcp

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/auth"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"

	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	mcpauth "github.com/thunder-id/thunderid/internal/system/mcp/auth"
	"github.com/thunder-id/thunderid/internal/system/security"
)

// Initialize initializes the MCP server and registers its routes with the provided mux.
func Initialize(
	mux *http.ServeMux,
	jwtService jwt.JWTServiceInterface,
) *mcpsdk.Server {
	cfg := config.GetServerRuntime().Config
	baseURL := config.GetServerURL(&cfg.Server)

	mcpURL := baseURL + MCPEndpointPath
	resourceMetadataURL := baseURL + OAuthProtectedResourceMetadataPath

	// Create MCP server and register standalone tools
	mcpServer := newServer()

	rootPerm := security.GetSystemRootPermission()

	tokenVerifier := mcpauth.NewTokenVerifier(jwtService, cfg.JWT.Issuer, mcpURL)
	httpHandler := mcpsdk.NewStreamableHTTPHandler(func(*http.Request) *mcpsdk.Server {
		return mcpServer
	}, nil)

	// Secure MCP handler with bearer token authentication
	securedHandler := auth.RequireBearerToken(tokenVerifier, &auth.RequireBearerTokenOptions{
		ResourceMetadataURL: resourceMetadataURL,
		Scopes:              []string{rootPerm},
	})(httpHandler)

	// Register protected resource metadata endpoint
	metadata := &oauthex.ProtectedResourceMetadata{
		Resource:             mcpURL,
		AuthorizationServers: []string{cfg.JWT.Issuer},
		ScopesSupported:      []string{rootPerm},
	}
	mux.Handle(OAuthProtectedResourceMetadataPath, auth.ProtectedResourceMetadataHandler(metadata))

	// Register MCP routes
	mux.Handle(MCPEndpointPath, securedHandler)
	mux.Handle(MCPEndpointPath+"/", securedHandler)

	return mcpServer
}
