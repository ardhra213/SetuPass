// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package email

import (
	"net/mail"
	"strings"
)

// IsValidEmail returns true for syntactically valid email addresses.
func IsValidEmail(emailAddr string) bool {
	// Reject CR/LF early (do not trim them away first) to prevent header injection.
	if strings.ContainsAny(emailAddr, "\r\n") {
		return false
	}

	emailAddr = strings.TrimSpace(emailAddr)
	if emailAddr == "" {
		return false
	}

	addr, err := mail.ParseAddress(emailAddr)
	if err != nil {
		return false
	}

	return addr.Address == emailAddr
}
