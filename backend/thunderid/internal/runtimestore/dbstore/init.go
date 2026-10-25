// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package dbstore provides a database-backed runtime store implementation.
package dbstore

import (
	"github.com/thunder-id/thunderid/internal/system/database/provider"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize creates and returns a new DBStore instance for the given deployment.
func Initialize(deploymentID string) (providers.RuntimeStoreProvider, providers.Transactioner, error) {
	dbProvider := provider.GetDBProvider()
	transactioner, error := dbProvider.GetRuntimeTransientDBTransactioner()
	if error != nil {
		return nil, nil, error
	}
	return newDBStore(dbProvider, deploymentID), transactioner, nil
}
