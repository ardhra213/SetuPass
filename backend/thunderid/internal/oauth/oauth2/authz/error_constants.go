// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authz

import "errors"

// errAuthorizationCodeNotFound is returned when an authorization code is not found in the database.
var errAuthorizationCodeNotFound = errors.New("authorization code not found")

// errAuthorizationCodeAlreadyConsumed is returned when an authorization code has already been consumed,
// indicating a potential replay attack.
var errAuthorizationCodeAlreadyConsumed = errors.New("authorization code already consumed")

// errAuthRequestNotFound is returned when an authorization request context is not found in the store.
var errAuthRequestNotFound = errors.New("authorization request context not found")

// errAssertionClaimInvalid is returned when a claim in the flow assertion has an unexpected shape
// (e.g. wrong JSON type). It distinguishes client-facing input errors from genuine internal decode failures.
var errAssertionClaimInvalid = errors.New("assertion claim is invalid")
