// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type ContextTestSuite struct {
	suite.Suite
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

func (suite *ContextTestSuite) TestWithKeyedTx() {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	suite.Require().NoError(err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	suite.Require().NoError(err)

	// Store transaction in context with key
	txCtx := WithKeyedTx(ctx, "test", tx)

	// Verify it's stored
	suite.NotNil(txCtx)
	suite.NotEqual(ctx, txCtx)
}

func (suite *ContextTestSuite) TestKeyedTxFromContext_WithTransaction() {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	suite.Require().NoError(err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	suite.Require().NoError(err)

	// Store transaction in context with key
	txCtx := WithKeyedTx(ctx, "test", tx)

	// Retrieve it
	retrievedTx := KeyedTxFromContext(txCtx, "test")
	suite.NotNil(retrievedTx)
	suite.Equal(tx, retrievedTx)
}

func (suite *ContextTestSuite) TestKeyedTxFromContext_WithoutTransaction() {
	ctx := context.Background()

	// Try to retrieve transaction from empty context
	tx := KeyedTxFromContext(ctx, "test")
	suite.Nil(tx)
}

func (suite *ContextTestSuite) TestKeyedTxFromContext_WrongType() {
	ctx := context.Background()

	// Store something other than *sql.Tx in the context
	ctx = context.WithValue(ctx, getTxContextKey("test"), "not a transaction")

	// Should return nil
	tx := KeyedTxFromContext(ctx, "test")
	suite.Nil(tx)
}

func (suite *ContextTestSuite) TestHasKeyedTx_WithTransaction() {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	suite.Require().NoError(err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	suite.Require().NoError(err)

	// Store transaction in context
	txCtx := WithKeyedTx(ctx, "test", tx)

	// Check if it has transaction
	suite.True(HasKeyedTx(txCtx, "test"))
}

func (suite *ContextTestSuite) TestHasKeyedTx_WithoutTransaction() {
	ctx := context.Background()

	// Check if empty context has transaction
	suite.False(HasKeyedTx(ctx, "test"))
}

func (suite *ContextTestSuite) TestHasKeyedTx_WrongType() {
	ctx := context.Background()

	// Store something other than *sql.Tx in the context
	ctx = context.WithValue(ctx, getTxContextKey("test"), "not a transaction")

	// Should return false
	suite.False(HasKeyedTx(ctx, "test"))
}
