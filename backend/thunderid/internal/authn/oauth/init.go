// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package oauth

import (
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/idp"
	syshttp "github.com/thunder-id/thunderid/internal/system/http"
)

// Initialize initializes the OAuth authentication service.
func Initialize(idpSvc idp.IDPServiceInterface,
	entityProvider entityprovider.EntityProviderInterface) OAuthAuthnServiceInterface {
	httpClient := syshttp.NewHTTPClient()
	return newOAuthAuthnService(httpClient, idpSvc, entityProvider)
}
