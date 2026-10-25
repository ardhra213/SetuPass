// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocationcache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/tests/mocks/database/providermock"
)

const testDeploymentID = "test-deployment-id"

type DBSourceTestSuite struct {
	suite.Suite
	mockDBProvider *providermock.DBProviderInterfaceMock
	mockDBClient   *providermock.DBClientInterfaceMock
	source         *dbSource
}

func TestDBSourceTestSuite(t *testing.T) {
	suite.Run(t, new(DBSourceTestSuite))
}

func (suite *DBSourceTestSuite) SetupTest() {
	suite.mockDBProvider = providermock.NewDBProviderInterfaceMock(suite.T())
	suite.mockDBClient = providermock.NewDBClientInterfaceMock(suite.T())
	suite.source = &dbSource{
		dbProvider:   suite.mockDBProvider,
		deploymentID: testDeploymentID,
	}
}

func (suite *DBSourceTestSuite) TestSnapshot_Success() {
	expiry := time.Now().Add(time.Hour).UTC()
	revokedAt := time.Now().UTC()
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{
			{"jti": "jti-1", "expiry_time": expiry},
			{"jti": "jti-2", "expiry_time": expiry},
		}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokenFamilies,
		criterionTypeTokenFamily, mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{
			{"criterion_value": "tfid-1", "expiry_time": expiry},
		}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedSubjects,
		criterionTypeSubject, mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{
			{"criterion_value": "user-1", "reason": "role_assignment_removed",
				"revoked_at": revokedAt, "expiry_time": expiry},
		}, nil)

	snapshot, err := suite.source.Snapshot(context.Background())

	suite.Require().NoError(err)
	assert.Len(suite.T(), snapshot.Tokens, 2)
	assert.Equal(suite.T(), "jti-1", snapshot.Tokens[0].Value)
	assert.Equal(suite.T(), expiry, snapshot.Tokens[0].ExpiryTime)
	assert.Len(suite.T(), snapshot.Families, 1)
	assert.Equal(suite.T(), "tfid-1", snapshot.Families[0].Value)
	assert.Len(suite.T(), snapshot.Subjects, 1)
	assert.Equal(suite.T(), "user-1", snapshot.Subjects[0].Value)
	assert.Equal(suite.T(), revokedAt, snapshot.Subjects[0].RevokedAt)
	assert.True(suite.T(), snapshot.Subjects[0].Boundary)
}

func (suite *DBSourceTestSuite) TestSnapshot_Empty() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokenFamilies,
		criterionTypeTokenFamily, mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedSubjects,
		criterionTypeSubject, mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{}, nil)

	snapshot, err := suite.source.Snapshot(context.Background())

	suite.Require().NoError(err)
	assert.Empty(suite.T(), snapshot.Tokens)
	assert.Empty(suite.T(), snapshot.Families)
	assert.Empty(suite.T(), snapshot.Subjects)
}

func (suite *DBSourceTestSuite) TestSnapshot_SubjectQueryError() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).Return([]map[string]interface{}{}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokenFamilies,
		criterionTypeTokenFamily, mock.Anything, testDeploymentID).Return([]map[string]interface{}{}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedSubjects,
		criterionTypeSubject, mock.Anything, testDeploymentID).Return(nil, errors.New("query error"))

	snapshot, err := suite.source.Snapshot(context.Background())

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), snapshot.Subjects)
	assert.Contains(suite.T(), err.Error(), "error reading revoked subject snapshot")
}

func (suite *DBSourceTestSuite) TestSnapshot_DBClientError() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(nil, errors.New("db client error"))

	snapshot, err := suite.source.Snapshot(context.Background())

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), snapshot.Tokens)
	assert.Contains(suite.T(), err.Error(), "db client error")
}

func (suite *DBSourceTestSuite) TestSnapshot_QueryError() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).
		Return(nil, errors.New("query error"))

	snapshot, err := suite.source.Snapshot(context.Background())

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), snapshot.Tokens)
	assert.Contains(suite.T(), err.Error(), "error reading revoked token snapshot")
}

func (suite *DBSourceTestSuite) TestSnapshot_TokenFamilyQueryError() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{}, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokenFamilies,
		criterionTypeTokenFamily, mock.Anything, testDeploymentID).
		Return(nil, errors.New("query error"))

	snapshot, err := suite.source.Snapshot(context.Background())

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), snapshot.Families)
	assert.Contains(suite.T(), err.Error(), "error reading revoked token family snapshot")
}

func (suite *DBSourceTestSuite) TestSnapshot_InvalidJTI() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{
			{"jti": "", "expiry_time": time.Now().Add(time.Hour)},
		}, nil)

	snapshot, err := suite.source.Snapshot(context.Background())

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), snapshot.Tokens)
	assert.Contains(suite.T(), err.Error(), "jti")
}

func (suite *DBSourceTestSuite) TestSnapshot_InvalidExpiryTime() {
	suite.mockDBProvider.On("GetRuntimePersistentDBClient").Return(suite.mockDBClient, nil)
	suite.mockDBClient.On("QueryContext", mock.Anything, querySnapshotRevokedTokens,
		mock.Anything, testDeploymentID).
		Return([]map[string]interface{}{
			{"jti": "jti-1", "expiry_time": 12345},
		}, nil)

	snapshot, err := suite.source.Snapshot(context.Background())

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), snapshot.Tokens)
	assert.Contains(suite.T(), err.Error(), "error parsing revocation snapshot")
}
