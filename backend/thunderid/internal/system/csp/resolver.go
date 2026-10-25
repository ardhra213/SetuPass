// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package csp

import (
	"context"
	"sync"

	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// ConfigReader reads the merged (effective) value of a server-config section. The server-config
// service satisfies it; this narrow interface keeps this package from importing serverconfig.
type ConfigReader interface {
	GetMergedConfig(ctx context.Context, name string) (any, *common.ServiceError)
}

var (
	readerMu sync.RWMutex
	reader   ConfigReader
)

// InitializeConfigReader installs the server-config reader the resolver reads the csp section from.
// It is called once at the composition root.
func InitializeConfigReader(r ConfigReader) {
	readerMu.Lock()
	defer readerMu.Unlock()
	reader = r
}

// Resolve returns the effective csp policy. When no reader is installed or the section cannot be read,
// it falls back to the zero PolicyConfig, which is the deny-first baseline in report-only mode, so the
// server always emits a safe policy.
func Resolve(ctx context.Context) PolicyConfig {
	readerMu.RLock()
	r := reader
	readerMu.RUnlock()
	if r == nil {
		return PolicyConfig{}
	}

	value, svcErr := r.GetMergedConfig(ctx, configSectionCSP)
	if svcErr != nil {
		log.GetLogger().With(log.String(log.LoggerKeyComponentName, "CSP")).Warn(ctx,
			"Failed to read the csp server config; falling back to the deny-first report-only default",
			log.String("code", svcErr.Code))
		return PolicyConfig{}
	}
	cfg, ok := value.(PolicyConfig)
	if !ok {
		return PolicyConfig{}
	}
	return cfg
}
