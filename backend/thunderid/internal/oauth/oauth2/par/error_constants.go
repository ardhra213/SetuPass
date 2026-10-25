// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package par implements OAuth 2.0 Pushed Authorization Requests (RFC 9126).
package par

import "errors"

var errInvalidRequestURI = errors.New("invalid request_uri format")

var errRequestURINotFound = errors.New("request_uri not found, expired, or already consumed")

// ErrPARResolutionFailed indicates a server-side failure while resolving a PAR request.
var ErrPARResolutionFailed = errors.New("failed to resolve pushed authorization request")

var errClientIDMismatch = errors.New("client_id does not match the pushed authorization request")
