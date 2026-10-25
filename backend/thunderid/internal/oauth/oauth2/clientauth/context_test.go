// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package clientauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

type ContextTestSuite struct {
	suite.Suite
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

func (suite *ContextTestSuite) TestGetOAuthClient_WithNilContext() {
	client := GetOAuthClient(context.Background())
	assert.Nil(suite.T(), client)
}

func (suite *ContextTestSuite) TestGetOAuthClient_WithEmptyContext() {
	ctx := context.Background()
	client := GetOAuthClient(ctx)
	assert.Nil(suite.T(), client)
}

func (suite *ContextTestSuite) TestGetOAuthClient_WithExistingClient() {
	expectedClient := &OAuthClientInfo{
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		OAuthApp: &providers.OAuthClient{
			ClientID: "test-client-id",
		},
	}

	ctx := withOAuthClient(context.Background(), expectedClient)
	client := GetOAuthClient(ctx)

	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), expectedClient.ClientID, client.ClientID)
	assert.Equal(suite.T(), expectedClient.ClientSecret, client.ClientSecret)
	assert.NotNil(suite.T(), client.OAuthApp)
}

func (suite *ContextTestSuite) TestWithOAuthClient() {
	expectedClient := &OAuthClientInfo{
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		OAuthApp: &providers.OAuthClient{
			ClientID:                "test-client-id",
			TokenEndpointAuthMethod: providers.TokenEndpointAuthMethodClientSecretPost,
		},
	}

	ctx := withOAuthClient(context.Background(), expectedClient)
	client := GetOAuthClient(ctx)

	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), expectedClient.ClientID, client.ClientID)
	assert.Equal(suite.T(), expectedClient.OAuthApp.ClientID, client.OAuthApp.ClientID)
}

func (suite *ContextTestSuite) TestWithOAuthClient_NilContext() {
	expectedClient := &OAuthClientInfo{
		ClientID: "test-client-id",
	}

	ctx := withOAuthClient(context.Background(), expectedClient)
	client := GetOAuthClient(ctx)

	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), expectedClient.ClientID, client.ClientID)
}

func (suite *ContextTestSuite) TestGetOAuthClient_WithWrongType() {
	ctx := context.WithValue(context.Background(), OAuthClientKey, "wrong-type")
	client := GetOAuthClient(ctx)
	assert.Nil(suite.T(), client)
}

func (suite *ContextTestSuite) TestGetOAuthClient_WithNilValue() {
	ctx := context.WithValue(context.Background(), OAuthClientKey, nil)
	client := GetOAuthClient(ctx)
	assert.Nil(suite.T(), client)
}

func (suite *ContextTestSuite) TestWithOAuthClient_NilClient() {
	ctx := withOAuthClient(context.Background(), nil)
	client := GetOAuthClient(ctx)
	assert.Nil(suite.T(), client)
}

func (suite *ContextTestSuite) TestWithOAuthClient_ContextChaining() {
	client1 := &OAuthClientInfo{
		ClientID: "client-1",
	}
	client2 := &OAuthClientInfo{
		ClientID: "client-2",
	}

	ctx1 := withOAuthClient(context.Background(), client1)
	ctx2 := withOAuthClient(ctx1, client2)

	client := GetOAuthClient(ctx2)
	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), "client-2", client.ClientID)
}

func (suite *ContextTestSuite) TestGetOAuthClient_WithEmptyClientInfo() {
	clientInfo := &OAuthClientInfo{
		ClientID:     "",
		ClientSecret: "",
		OAuthApp:     nil,
	}

	ctx := withOAuthClient(context.Background(), clientInfo)
	client := GetOAuthClient(ctx)

	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), "", client.ClientID)
}
