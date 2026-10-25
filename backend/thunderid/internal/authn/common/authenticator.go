// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"
	"sync"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

var (
	authenticatorRegistry map[string]AuthenticatorMeta
	registryMu            sync.RWMutex
)

func init() {
	authenticatorRegistry = make(map[string]AuthenticatorMeta)
}

// RegisterAuthenticator registers an authenticator's metadata in the registry.
// Each authenticator service should call this during its initialization.
func RegisterAuthenticator(meta AuthenticatorMeta) {
	registryMu.Lock()
	defer registryMu.Unlock()
	authenticatorRegistry[meta.Name] = meta
}

// getAuthenticatorMetaData returns the authenticator metadata for the given authenticator.
func getAuthenticatorMetaData(name string) *AuthenticatorMeta {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if auth, ok := authenticatorRegistry[name]; ok {
		return &auth
	}
	return nil
}

// GetAuthenticatorFactors returns the authentication factors for the given authenticator.
func GetAuthenticatorFactors(name string) []AuthenticationFactor {
	if auth := getAuthenticatorMetaData(name); auth != nil {
		return auth.Factors
	}
	return []AuthenticationFactor{}
}

// GetAuthenticatorNameForIDPType returns the authenticator name for a given IDP type.
func GetAuthenticatorNameForIDPType(idpType providers.IDPType) (string, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if idpType != "" {
		for _, meta := range authenticatorRegistry {
			if meta.AssociatedIDP == idpType {
				return meta.Name, nil
			}
		}
	}

	return "", errors.New("no authenticator found for the given IDP type")
}
