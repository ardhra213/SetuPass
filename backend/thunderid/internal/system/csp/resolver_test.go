// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package csp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

type stubReader struct {
	value any
	err   *common.ServiceError
}

func (s stubReader) GetMergedConfig(_ context.Context, _ string) (any, *common.ServiceError) {
	return s.value, s.err
}

func TestResolve(t *testing.T) {
	t.Cleanup(func() { InitializeConfigReader(nil) })

	t.Run("nil reader falls back to the default report-only policy", func(t *testing.T) {
		InitializeConfigReader(nil)
		assert.Equal(t, PolicyConfig{}, Resolve(context.Background()))
	})

	t.Run("returns the reader's merged policy", func(t *testing.T) {
		want := PolicyConfig{ReportOnly: boolPtr(false), ReportURI: "/r"}
		InitializeConfigReader(stubReader{value: want})
		assert.Equal(t, want, Resolve(context.Background()))
	})

	t.Run("service error falls back to the default", func(t *testing.T) {
		InitializeConfigReader(stubReader{err: &common.InternalServerError})
		assert.Equal(t, PolicyConfig{}, Resolve(context.Background()))
	})

	t.Run("unexpected value type falls back to the default", func(t *testing.T) {
		InitializeConfigReader(stubReader{value: "not a policy"})
		assert.Equal(t, PolicyConfig{}, Resolve(context.Background()))
	})
}
