// SPDX-License-Identifier: Apache-2.0
package providers

import "context"

const DeviceSecretTokenType = "urn:openid:params:token-type:device-secret"

// DeviceSecretContext carries server-authorized state, never caller-supplied identity.
type DeviceSecretContext struct {
	Subject          string                 `json:"sub"`
	ClientID         string                 `json:"client_id"`
	AuthID           string                 `json:"auth_id"`
	SID              string                 `json:"sid"`
	DPoPJkt          string                 `json:"jkt"`
	Scopes           []string               `json:"scope"`
	ACR              string                 `json:"acr"`
	AuthTime         int64                  `json:"auth_time"`
	ExpiresAt        int64                  `json:"exp"`
	AttributeCacheID string                 `json:"aci"`
	TokenFamilyID    string                 `json:"tfid"`
	Claims           map[string]interface{} `json:"claims"`
}

type DeviceSecretResult struct{ Secret, Hash, SID string }

// DeviceSecretProvider is supplied by the embedding application through its ActorProvider.
// Validate must authenticate both tokens, their binding, the live session, proof key and client family.
type DeviceSecretProvider interface {
	GenerateDeviceSecret(context.Context, DeviceSecretContext) (*DeviceSecretResult, error)
	ValidateDeviceSecret(ctx context.Context, secret, idToken, childClientID, proofJkt string, scopes []string) (*DeviceSecretContext, error)
}
