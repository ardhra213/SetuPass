// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package corstest provides helpers for installing test CORS matchers.
package corstest

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/thunder-id/thunderid/internal/system/cors"
	"github.com/thunder-id/thunderid/tests/mocks/corsmock"
)

// InstallMatcherEntries installs a CORS matcher for the duration of the test, serving the given entries as
// the writable layer over an empty read-only layer.
func InstallMatcherEntries(t *testing.T, entries cors.OriginEntries) {
	reader := corsmock.NewServerConfigReaderMock(t)
	reader.EXPECT().
		GetReadOnlyConfig(mock.Anything, "cors").
		Return(cors.OriginConfig{AllowedOrigins: cors.OriginEntries{}}, nil).
		Maybe()
	reader.EXPECT().
		GetWritableConfig(mock.Anything, "cors").
		Return(cors.OriginConfig{AllowedOrigins: entries}, nil).
		Maybe()
	cors.InitializeDynamicMatcher(reader)
	t.Cleanup(func() { cors.InitializeDynamicMatcher(nil) })
}
