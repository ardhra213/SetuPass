// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	authnprovidercm "github.com/thunder-id/thunderid/internal/authnprovider/common"
)

// Field keys for entity system attributes.
const (
	fieldName         = "name"
	fieldDescription  = "description"
	fieldClientID     = "clientId"
	fieldClientSecret = authnprovidercm.CredentialTypeClientSecret
	fieldFlowSecret   = authnprovidercm.CredentialTypeFlowSecret
)

// Field keys for application config properties.
const (
	propURL         = "url"
	propLogoURL     = "logo_url"
	propTosURI      = "tos_uri"
	propPolicyURI   = "policy_uri"
	propContacts    = "contacts"
	propType        = "type"
	propTemplate    = "template"
	propMetadata    = "metadata"
	propOAuthConfig = "oauth_config"
)
