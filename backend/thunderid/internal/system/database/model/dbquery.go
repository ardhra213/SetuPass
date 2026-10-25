// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

// DBQueryInterface defines the interface for database queries.
type DBQueryInterface interface {
	GetID() string
	GetQuery(dbType string) string
}

var _ DBQueryInterface = (*DBQuery)(nil)

// DBQuery represents database queries with an identifier and the SQL query string.
type DBQuery struct {
	// ID is the unique identifier for the query.
	ID            string `json:"id"`
	Query         string `json:"query"`
	PostgresQuery string `json:"postgres_query,omitempty"`
	SQLiteQuery   string `json:"sqlite_query,omitempty"`
}

// GetID returns the unique identifier for the query.
func (d *DBQuery) GetID() string {
	return d.ID
}

// GetQuery returns the appropriate query for the specified database type.
func (d *DBQuery) GetQuery(dbType string) string {
	switch dbType {
	case "postgres":
		if d.PostgresQuery != "" {
			return d.PostgresQuery
		}
	case "sqlite":
		if d.SQLiteQuery != "" {
			return d.SQLiteQuery
		}
	}
	// Fall back to the default query
	return d.Query
}
