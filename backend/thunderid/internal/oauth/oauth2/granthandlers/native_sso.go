// SPDX-License-Identifier: Apache-2.0
package granthandlers

import (
	"context"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/dpop"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/model"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/tokenservice"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
	"strings"
	"time"
)

func (h *tokenExchangeGrantHandler) handleNativeSSO(ctx context.Context, r *model.TokenRequest, app *providers.OAuthClient) (*model.TokenResponseDTO, *model.ErrorResponse) {
	deny := func(code string) (*model.TokenResponseDTO, *model.ErrorResponse) {
		return nil, &model.ErrorResponse{Error: code, ErrorDescription: "Native SSO request denied; sign in to Setu ID again"}
	}
	if app.DeviceSecretProvider == nil || r.SubjectTokenType != "urn:ietf:params:oauth:token-type:id_token" || dpop.GetJkt(ctx) == "" || len(r.Resources) > 0 || len(r.Audiences) > 0 || (r.RequestedTokenType != "" && r.RequestedTokenType != "urn:ietf:params:oauth:token-type:access_token") {
		return deny("invalid_request")
	}
	scopes := strings.Fields(r.Scope)
	if len(scopes) == 0 {
		return deny("invalid_scope")
	}
	state, err := app.DeviceSecretProvider.ValidateDeviceSecret(ctx, r.ActorToken, r.SubjectToken, r.ClientID, dpop.GetJkt(ctx), scopes)
	if err != nil {
		return deny("invalid_grant")
	}
	ttl := min(int64(120), state.ExpiresAt-time.Now().Unix())
	if ttl <= 0 {
		return deny("invalid_grant")
	}
	at, err := h.tokenBuilder.BuildAccessToken(ctx, &tokenservice.AccessTokenBuildContext{Subject: state.Subject, Audiences: []string{app.ResolveDefaultAudience(r.ClientID)}, ClientID: r.ClientID, Scopes: scopes, AttributeCacheID: state.AttributeCacheID, GrantType: string(providers.GrantTypeTokenExchange), OAuthApp: app, DPoPJkt: dpop.GetJkt(ctx), TokenFamilyID: state.TokenFamilyID, ValidityPeriod: ttl})
	if err != nil {
		return deny("server_error")
	}
	response := &model.TokenResponseDTO{AccessToken: *at}
	for _, s := range scopes {
		if s == "openid" {
			id, err := h.tokenBuilder.BuildIDToken(ctx, &tokenservice.IDTokenBuildContext{Subject: state.Subject, Audience: r.ClientID, Scopes: scopes, OAuthApp: app, AuthTime: state.AuthTime, CompletedACR: state.ACR, SessionID: state.SID, ValidityPeriod: ttl})
			if err != nil {
				return deny("server_error")
			}
			response.IDToken = *id
		}
	}
	return response, nil
}
