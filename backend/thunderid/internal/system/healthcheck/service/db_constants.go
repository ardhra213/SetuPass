// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package service provides health check-related business logic and operations.
package service

import "github.com/thunder-id/thunderid/internal/system/database/model"

var queryConfigDBTable = model.DBQuery{
	ID:    "HLC-00001",
	Query: "SELECT 1",
}

var queryRuntimeTransientDBTable = model.DBQuery{
	ID:    "HLC-00002",
	Query: "SELECT 1",
}

var queryEntityDBTable = model.DBQuery{
	ID:    "HLC-00003",
	Query: "SELECT 1",
}
