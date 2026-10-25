// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package ciba

import "errors"

// ErrCIBARequestNotFound is returned when a CIBA authentication request is not found in the store.
var ErrCIBARequestNotFound = errors.New("ciba authentication request not found")
