// SPDX-License-Identifier: Apache-2.0
package granthandlers

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/dpop"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/model"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/tokenservice"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
	"github.com/thunder-id/thunderid/tests/mocks/oauth/oauth2/tokenservicemock"
	"testing"
	"time"
)

type testDeviceProvider struct {
	providers.DeviceSecretProvider
	deny bool
}

func (p testDeviceProvider) ValidateDeviceSecret(_ context.Context, secret, id, child, jkt string, scopes []string) (*providers.DeviceSecretContext, error) {
	if p.deny || secret != "opaque" || id != "parent-id" || child != "child" || jkt != "key" {
		return nil, errors.New("denied")
	}
	return &providers.DeviceSecretContext{Subject: "person", SID: "session", AuthTime: 123, ACR: "otp", ExpiresAt: time.Now().Unix() + 90, AttributeCacheID: "attributes", TokenFamilyID: "family"}, nil
}
func TestNativeExchangeBuildsChildTokens(t *testing.T) {
	b := tokenservicemock.NewTokenBuilderInterfaceMock(t)
	b.On("BuildAccessToken", mock.Anything, mock.MatchedBy(func(c *tokenservice.AccessTokenBuildContext) bool {
		return c.Subject == "person" && c.ClientID == "child" && c.DPoPJkt == "key" && c.AttributeCacheID == "attributes" && c.TokenFamilyID == "family" && c.ValidityPeriod > 0 && c.ValidityPeriod <= 90
	})).Return(&model.TokenDTO{Token: "child-at", TokenType: "DPoP", Scopes: []string{"openid", "email"}, ExpiresIn: 90}, nil).Once()
	b.On("BuildIDToken", mock.Anything, mock.MatchedBy(func(c *tokenservice.IDTokenBuildContext) bool {
		return c.Audience == "child" && c.Subject == "person" && c.SessionID == "session" && c.DeviceSecret == "" && c.ValidityPeriod <= 90
	})).Return(&model.TokenDTO{Token: "child-id"}, nil).Once()
	h := &tokenExchangeGrantHandler{tokenBuilder: b}
	r := &model.TokenRequest{GrantType: string(providers.GrantTypeTokenExchange), ClientID: "child", Scope: "openid email", ActorTokenType: providers.DeviceSecretTokenType, ActorToken: "opaque", SubjectTokenType: "urn:ietf:params:oauth:token-type:id_token", SubjectToken: "parent-id"}
	app := &providers.OAuthClient{DeviceSecretProvider: testDeviceProvider{}}
	require.Nil(t, h.ValidateGrant(context.Background(), r, app))
	out, e := h.HandleGrant(dpop.WithJkt(context.Background(), "key"), r, app)
	require.Nil(t, e)
	require.Equal(t, "child-at", out.AccessToken.Token)
	require.Equal(t, "child-id", out.IDToken.Token)
	require.Empty(t, out.DeviceSecret)
	app.DeviceSecretProvider = testDeviceProvider{deny: true}
	_, e = h.HandleGrant(dpop.WithJkt(context.Background(), "key"), r, app)
	require.Equal(t, "invalid_grant", e.Error)
	_, e = h.HandleGrant(context.Background(), r, app)
	require.Equal(t, "invalid_request", e.Error)
	r.Resources = []string{"https://untrusted.example"}
	_, e = h.HandleGrant(dpop.WithJkt(context.Background(), "key"), r, app)
	require.Equal(t, "invalid_request", e.Error)
}
