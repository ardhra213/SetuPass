// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/database/model"
)

type DBProviderTestSuite struct {
	suite.Suite
	mockDB sqlmock.Sqlmock
}

func TestDBProviderTestSuite(t *testing.T) {
	suite.Run(t, new(DBProviderTestSuite))
}

func (suite *DBProviderTestSuite) SetupTest() {
	_, mock, err := sqlmock.New()
	suite.Require().NoError(err)
	suite.mockDB = mock

	// Reset global config before each test
	config.ResetServerRuntime()

	// Initialize a dummy config
	dummyConfig := &config.Config{
		Database: config.DatabaseConfig{
			Config: config.DataSource{
				Type: "postgres", Postgres: config.PostgresDataSource{Name: "identity"},
			},
			RuntimeTransient: config.DataSource{
				Type: "postgres", Postgres: config.PostgresDataSource{Name: "runtime"},
			},
			Entity: config.DataSource{
				Type: "postgres", Postgres: config.PostgresDataSource{Name: "entity"},
			},
			RuntimePersistent: config.DataSource{
				Type: "postgres", Postgres: config.PostgresDataSource{Name: "operation"},
			},
		},
	}
	err = config.InitializeServerRuntime(".", dummyConfig)
	suite.Require().NoError(err)
}

func (suite *DBProviderTestSuite) TearDownTest() {
	config.ResetServerRuntime()
}

func (suite *DBProviderTestSuite) TestGetEntityDBTransactioner_Success() {
	// Create a mock DB connection
	db, _, err := sqlmock.New()
	suite.Require().NoError(err)
	defer func() {
		_ = db.Close()
	}()

	// Manually construct the provider with an initialized client
	provider := &dbProvider{
		entityClient: NewDBClient(model.NewDB(db), "postgres", "entity", retryConfig{}),
	}

	// Test getting the transactioner
	txer, err := provider.GetEntityDBTransactioner()
	suite.NoError(err)
	suite.NotNil(txer)
}

func (suite *DBProviderTestSuite) TestGetRuntimeTransientDBTransactioner_Success() {
	// Create a mock DB connection
	db, _, err := sqlmock.New()
	suite.Require().NoError(err)
	defer func() {
		_ = db.Close()
	}()

	// Manually construct the provider with an initialized client
	provider := &dbProvider{
		runtimeTransientClient: NewDBClient(model.NewDB(db), "postgres", "runtime", retryConfig{}),
	}

	// Test getting the transactioner
	txer, err := provider.GetRuntimeTransientDBTransactioner()
	suite.NoError(err)
	suite.NotNil(txer)
}

func (suite *DBProviderTestSuite) TestGetRuntimePersistentDBTransactioner_Success() {
	// Create a mock DB connection
	db, _, err := sqlmock.New()
	suite.Require().NoError(err)
	defer func() {
		_ = db.Close()
	}()

	// Manually construct the provider with an initialized client
	provider := &dbProvider{
		runtimePersistentClient: NewDBClient(model.NewDB(db), "postgres", "operation", retryConfig{}),
	}

	// Test getting the transactioner
	txer, err := provider.GetRuntimePersistentDBTransactioner()
	suite.NoError(err)
	suite.NotNil(txer)
}
