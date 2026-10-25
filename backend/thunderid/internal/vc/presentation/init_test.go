// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package presentation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
)

func setupDefinitionConfig(t *testing.T, store string, declarativeEnabled bool) {
	t.Helper()
	config.ResetServerRuntime()
	t.Cleanup(config.ResetServerRuntime)
	require.NoError(t, config.InitializeServerRuntime("", &config.Config{
		OpenID4VP:            config.OpenID4VPConfig{Store: store},
		DeclarativeResources: config.DeclarativeResources{Enabled: declarativeEnabled},
	}))
}

func TestGetDefinitionStoreMode_ExplicitMutable(t *testing.T) {
	setupDefinitionConfig(t, "mutable", false)
	mode, err := getDefinitionStoreMode()
	require.NoError(t, err)
	assert.Equal(t, serverconst.StoreModeMutable, mode)
}

func TestGetDefinitionStoreMode_ExplicitDeclarative(t *testing.T) {
	setupDefinitionConfig(t, "declarative", false)
	mode, err := getDefinitionStoreMode()
	require.NoError(t, err)
	assert.Equal(t, serverconst.StoreModeDeclarative, mode)
}

func TestGetDefinitionStoreMode_ExplicitComposite(t *testing.T) {
	setupDefinitionConfig(t, "composite", false)
	mode, err := getDefinitionStoreMode()
	require.NoError(t, err)
	assert.Equal(t, serverconst.StoreModeComposite, mode)
}

func TestGetDefinitionStoreMode_CaseInsensitive(t *testing.T) {
	setupDefinitionConfig(t, "  Composite  ", false)
	mode, err := getDefinitionStoreMode()
	require.NoError(t, err)
	assert.Equal(t, serverconst.StoreModeComposite, mode)
}

func TestGetDefinitionStoreMode_InvalidIsError(t *testing.T) {
	setupDefinitionConfig(t, "bogus", true)
	_, err := getDefinitionStoreMode()
	assert.Error(t, err)
}

func TestGetDefinitionStoreMode_FallbackDeclarativeEnabled(t *testing.T) {
	setupDefinitionConfig(t, "", true)
	mode, err := getDefinitionStoreMode()
	require.NoError(t, err)
	assert.Equal(t, serverconst.StoreModeDeclarative, mode)
}

func TestGetDefinitionStoreMode_FallbackDeclarativeDisabled(t *testing.T) {
	setupDefinitionConfig(t, "", false)
	mode, err := getDefinitionStoreMode()
	require.NoError(t, err)
	assert.Equal(t, serverconst.StoreModeMutable, mode)
}

func TestRegisterRoutesRegistersEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	handler := newDefinitionHandler(NewPresentationDefinitionServiceInterfaceMock(t))
	registerRoutes(mux, handler)

	// The OPTIONS preflight handlers respond without invoking the service.
	for _, target := range []string{definitionsPath, definitionsPath + "/some-id"} {
		req := httptest.NewRequest(http.MethodOptions, target, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
