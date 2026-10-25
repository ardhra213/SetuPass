// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package jwksresolver

import (
	syshttp "github.com/thunder-id/thunderid/internal/system/http"
)

// Initialize creates and returns a new Resolver configured with the given HTTP client.
func Initialize(httpClient syshttp.HTTPClientInterface) *Resolver {
	return newJWKSResolver(httpClient)
}
