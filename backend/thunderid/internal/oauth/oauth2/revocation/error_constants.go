// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocation

import "errors"

// ErrTokenRevoked indicates the presented token's JTI is on the deny list.
var ErrTokenRevoked = errors.New("token has been revoked")

// ErrEnforcementUnavailable indicates the deny list could not be consulted (runtime persistent DB
// unavailable or the circuit is open). Under the fail-closed policy callers MUST reject the token.
var ErrEnforcementUnavailable = errors.New("token revocation enforcement is unavailable")
