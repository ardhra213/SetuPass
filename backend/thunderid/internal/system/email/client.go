// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package email provides a reusable email sending capability for the server.
// It defines a common interface and implementations for sending emails via various transports.
package email

import "context"

// EmailClientInterface defines the interface for sending emails.
type EmailClientInterface interface {
	// Send sends an email using the provided EmailData.
	Send(ctx context.Context, emailData EmailData) error
}
