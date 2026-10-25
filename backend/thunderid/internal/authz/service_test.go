// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authz

import (
	"context"
	"errors"
	"testing"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/authz/engine"
	enginemock "github.com/thunder-id/thunderid/tests/mocks/authz/engine"
)

type AuthorizationServiceTestSuite struct {
	suite.Suite
	mockEngine *enginemock.AuthorizationEngineMock
	service    providers.AuthorizationProvider
}

func TestAuthorizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorizationServiceTestSuite))
}

func (suite *AuthorizationServiceTestSuite) SetupTest() {
	suite.mockEngine = enginemock.NewAuthorizationEngineMock(suite.T())
	suite.service = newAuthorizationService(suite.mockEngine)
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessSuccess() {
	request := providers.AccessEvaluationRequest{
		Subject:        providers.Subject{ID: "user1", GroupIDs: []string{"group1"}},
		ResourceServer: providers.AccessEvaluationResourceServer{ID: "document"},
		Permission:     providers.Permission{Name: "read"},
	}

	suite.mockEngine.On("EvaluateAccessBatch", mock.Anything,
		mock.MatchedBy(func(req engine.AccessEvaluationsRequest) bool {
			return len(req.Evaluations) == 1 &&
				req.Evaluations[0].Subject.ID == "user1" &&
				req.Evaluations[0].Subject.GroupIDs[0] == "group1" &&
				req.Evaluations[0].ResourceServer.ID == "document" &&
				req.Evaluations[0].Permission.Name == "read"
		})).
		Return(&engine.AccessEvaluationsResponse{
			Evaluations: []engine.AccessEvaluationResponse{{Decision: true}},
		}, nil)

	response, err := suite.service.EvaluateAccess(context.Background(), request)

	suite.Nil(err)
	suite.NotNil(response)
	suite.True(response.Decision)
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessBatchSuccess() {
	request := providers.AccessEvaluationsRequest{
		Evaluations: []providers.AccessEvaluationRequest{
			{
				Subject:        providers.Subject{ID: "user1", GroupIDs: []string{"group1"}},
				ResourceServer: providers.AccessEvaluationResourceServer{ID: "document"},
				Permission:     providers.Permission{Name: "read"},
			},
			{
				Subject:        providers.Subject{ID: "user1", GroupIDs: []string{"group1"}},
				ResourceServer: providers.AccessEvaluationResourceServer{ID: "document"},
				Permission:     providers.Permission{Name: "delete"},
			},
		},
	}

	suite.mockEngine.On("EvaluateAccessBatch", mock.Anything, mock.Anything).
		Return(&engine.AccessEvaluationsResponse{
			Evaluations: []engine.AccessEvaluationResponse{
				{Decision: true},
				{Decision: false},
			},
		}, nil)

	response, err := suite.service.EvaluateAccessBatch(context.Background(), request)

	suite.Nil(err)
	suite.NotNil(response)
	suite.Len(response.Evaluations, 2)
	suite.True(response.Evaluations[0].Decision)
	suite.False(response.Evaluations[1].Decision)
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessBatchReturnsContext() {
	request := providers.AccessEvaluationsRequest{
		Evaluations: []providers.AccessEvaluationRequest{
			{
				Subject:        providers.Subject{ID: "user1"},
				ResourceServer: providers.AccessEvaluationResourceServer{ID: "document"},
				Permission:     providers.Permission{Name: "read"},
			},
		},
	}
	decisionContext := map[string]interface{}{
		"reason": "requires_step_up",
	}

	suite.mockEngine.On("EvaluateAccessBatch", mock.Anything, mock.Anything).
		Return(&engine.AccessEvaluationsResponse{
			Evaluations: []engine.AccessEvaluationResponse{
				{Decision: false, Context: decisionContext},
			},
		}, nil)

	response, err := suite.service.EvaluateAccessBatch(context.Background(), request)

	suite.Nil(err)
	suite.NotNil(response)
	suite.Equal(decisionContext, response.Evaluations[0].Context)
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessPassesPropertiesToEngine() {
	subjectProperties := map[string]interface{}{"department": "Sales"}
	resourceProperties := map[string]interface{}{"owner": "user1"}
	actionProperties := map[string]interface{}{"method": "GET"}
	request := providers.AccessEvaluationRequest{
		Subject: providers.Subject{
			Type:       "user",
			ID:         "user1",
			Properties: subjectProperties,
		},
		ResourceServer: providers.AccessEvaluationResourceServer{
			ID:         "document",
			Properties: resourceProperties,
		},
		Permission: providers.Permission{
			Name:       "read",
			Properties: actionProperties,
		},
	}

	suite.mockEngine.On("EvaluateAccessBatch", mock.Anything,
		mock.MatchedBy(func(req engine.AccessEvaluationsRequest) bool {
			if len(req.Evaluations) != 1 {
				return false
			}
			evaluation := req.Evaluations[0]
			return suite.Equal(subjectProperties, evaluation.Subject.Properties) &&
				suite.Equal(resourceProperties, evaluation.ResourceServer.Properties) &&
				suite.Equal(actionProperties, evaluation.Permission.Properties)
		})).
		Return(&engine.AccessEvaluationsResponse{
			Evaluations: []engine.AccessEvaluationResponse{{Decision: true}},
		}, nil)

	response, err := suite.service.EvaluateAccess(context.Background(), request)

	suite.Nil(err)
	suite.NotNil(response)
	suite.True(response.Decision)
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessBatchEmpty() {
	response, err := suite.service.EvaluateAccessBatch(context.Background(), providers.AccessEvaluationsRequest{})

	suite.Nil(err)
	suite.NotNil(response)
	suite.Empty(response.Evaluations)
	suite.mockEngine.AssertNotCalled(suite.T(), "EvaluateAccessBatch")
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessBatchEngineError() {
	request := providers.AccessEvaluationsRequest{
		Evaluations: []providers.AccessEvaluationRequest{
			{
				Subject:        providers.Subject{ID: "user1"},
				ResourceServer: providers.AccessEvaluationResourceServer{ID: "document"},
				Permission:     providers.Permission{Name: "read"},
			},
		},
	}

	suite.mockEngine.On("EvaluateAccessBatch", mock.Anything, mock.Anything).
		Return((*engine.AccessEvaluationsResponse)(nil), errors.New("engine failed"))

	response, err := suite.service.EvaluateAccessBatch(context.Background(), request)

	suite.Nil(response)
	suite.NotNil(err)
	suite.Equal(tidcommon.InternalServerError.Code, err.Code)
}

func (suite *AuthorizationServiceTestSuite) TestEvaluateAccessEmptyEngineResponse() {
	request := providers.AccessEvaluationRequest{
		Subject:        providers.Subject{ID: "user1"},
		ResourceServer: providers.AccessEvaluationResourceServer{ID: "document"},
		Permission:     providers.Permission{Name: "read"},
	}

	suite.mockEngine.On("EvaluateAccessBatch", mock.Anything, mock.Anything).
		Return(&engine.AccessEvaluationsResponse{}, nil)

	response, err := suite.service.EvaluateAccess(context.Background(), request)

	suite.Nil(err)
	suite.NotNil(response)
	suite.False(response.Decision)
}
