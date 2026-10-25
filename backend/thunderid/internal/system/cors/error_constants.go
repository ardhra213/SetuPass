// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cors

import "errors"

// Sentinel errors returned by the parser and compiler. Callers may wrap these
// for additional context but should compare via errors.Is.
var (
	// ErrInvalidOrigin is returned by ParseOrigin when the request Origin
	// header is syntactically invalid as an HTTP(S) origin.
	ErrInvalidOrigin = errors.New("cors: invalid origin")

	// ErrInvalidLiteral is returned by Compile when a literal allowed-origin
	// entry cannot be parsed as a valid origin.
	ErrInvalidLiteral = errors.New("cors: invalid literal entry")

	// ErrInvalidRegex is returned by Compile when a regex allowed-origin
	// entry fails to compile under RE2.
	ErrInvalidRegex = errors.New("cors: invalid regex entry")

	// ErrEmptyEntry is returned by Compile when an entry carries no value
	// (empty literal string or empty regex pattern).
	ErrEmptyEntry = errors.New("cors: empty entry")

	// ErrWildcardLiteral is returned by Compile when an operator configures
	// the literal "*" entry. CORS does not allow combining the wildcard with
	// credentials, and the project does not support unauthenticated allow-all
	// either; operators must list explicit origins or use a regex entry.
	ErrWildcardLiteral = errors.New("cors: wildcard '*' literal is not supported")
)
