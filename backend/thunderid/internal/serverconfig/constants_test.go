// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package serverconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultResourceServerConfigNameSupported(t *testing.T) {
	assert.Equal(t, ConfigName("defaultResourceServer"), ConfigNameDefaultResourceServer)
	assert.True(t, ConfigNameDefaultResourceServer.IsValid())
	assert.Contains(t, supportedConfigNames, ConfigNameDefaultResourceServer)
}
