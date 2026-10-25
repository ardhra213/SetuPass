// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

// Package client defines the service and interfaces for sending messages.
package client

import (
	"context"
	"time"

	"github.com/thunder-id/thunderid/internal/notification/common"
)

// httpClientTimeout is the timeout duration for the HTTP client.
const httpClientTimeout = 10 * time.Second

// NotificationClientInterface defines the provider client interface for sending notifications.
type NotificationClientInterface interface {
	GetName() string
	IsChannelSupported(channel common.ChannelType) bool
	Send(ctx context.Context, channel common.ChannelType, data common.NotificationData) error
}
