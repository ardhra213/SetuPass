// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package revocationcache

import "errors"

// errTokenRevoked is returned by EnforcerInterface.EnsureNotRevoked when the token identifier is
// present in the cached deny list.
var errTokenRevoked = errors.New("token has been revoked")

// errUnsupportedSource is returned by Initialize when cfg.Source names a sync source that is not
// supported.
var errUnsupportedSource = errors.New("unsupported revocation sync source")
