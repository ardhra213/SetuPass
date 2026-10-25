// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package inmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	store := Initialize(testDeploymentID)

	ims, ok := store.(*inMemoryStore)
	require.True(t, ok)
	assert.Equal(t, testDeploymentID, ims.deploymentID)
	assert.NotNil(t, ims.data)
}
