// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package runtimestore provides the factory that selects a runtime store backend.
package runtimestore

import (
	"github.com/thunder-id/thunderid/internal/runtimestore/dbstore"
	"github.com/thunder-id/thunderid/internal/runtimestore/redisstore"
	dbprovider "github.com/thunder-id/thunderid/internal/system/database/provider"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize returns the runtime store provider backing the given runtime datasource type.
// Redis-backed runtimes use the Redis store; all others use the relational database store.
func Initialize(runtimeTransientDBType, deploymentID string) (
	providers.RuntimeStoreProvider, providers.Transactioner, error) {
	if runtimeTransientDBType == dbprovider.DataSourceTypeRedis {
		return redisstore.Initialize(deploymentID)
	}
	return dbstore.Initialize(deploymentID)
}
