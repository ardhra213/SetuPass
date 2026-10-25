// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package par

import (
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/model"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// pushedAuthorizationRequest holds the stored PAR data.
type pushedAuthorizationRequest struct {
	ClientID         string
	OAuthParameters  model.OAuthParameters
	InitiatorRequest providers.InitiatorRequest
}

// parResponse represents the PAR endpoint success response.
type parResponse struct {
	RequestURI string `json:"request_uri"`
	ExpiresIn  int64  `json:"expires_in"`
}
