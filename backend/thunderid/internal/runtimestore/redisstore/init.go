// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package redisstore provides a Redis-backed runtime store implementation.
package redisstore

import (
	"github.com/thunder-id/thunderid/internal/system/transaction"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Initialize creates and returns a new RedisStore instance for the given deployment.
func Initialize(deploymentID string) (providers.RuntimeStoreProvider, providers.Transactioner, error) {
	return newRedisStore(deploymentID), transaction.NewNoOpTransactioner(), nil
}
