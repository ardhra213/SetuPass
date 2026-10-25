// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package email

// EmailData represents the data structure for an email message.
type EmailData struct {
	To      []string `json:"to"`      // recipient email addresses
	CC      []string `json:"cc"`      // CC email addresses (optional)
	BCC     []string `json:"bcc"`     // BCC email addresses (optional)
	Subject string   `json:"subject"` // email subject
	Body    string   `json:"body"`    // email body content
	IsHTML  bool     `json:"is_html"` // true for HTML content, false for plain text
}

type smtpConfig struct {
	host                 string
	port                 int
	username             string
	password             string
	from                 string
	useTLS               bool
	enableAuthentication bool
}

type smtpClient struct {
	config smtpConfig
}
