// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package notification

import (
	"context"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"github.com/thunder-id/thunderid/internal/notification/client"
	"github.com/thunder-id/thunderid/internal/notification/common"
	"github.com/thunder-id/thunderid/internal/system/log"
)

// NotificationSenderServiceInterface defines the interface for sending notification messages.
type NotificationSenderServiceInterface interface {
	Send(ctx context.Context, channel common.ChannelType, senderID string,
		data common.NotificationData) *tidcommon.ServiceError
}

// notificationSenderService implements NotificationSenderServiceInterface.
type notificationSenderService struct {
	senderMgtService NotificationSenderMgtSvcInterface
	clientFactory    client.ClientFactoryInterface
	logger           *log.Logger
}

// newNotificationSenderService returns a new instance of NotificationSenderServiceInterface.
func newNotificationSenderService(
	senderMgtService NotificationSenderMgtSvcInterface,
	clientFactory client.ClientFactoryInterface) NotificationSenderServiceInterface {
	return &notificationSenderService{
		senderMgtService: senderMgtService,
		clientFactory:    clientFactory,
		logger:           log.GetLogger().With(log.String(log.LoggerKeyComponentName, "NotificationSenderService")),
	}
}

// Send looks up the sender by ID and dispatches the notification via the specified channel.
func (s *notificationSenderService) Send(ctx context.Context, channel common.ChannelType, senderID string,
	data common.NotificationData) *tidcommon.ServiceError {
	sender, svcErr := s.senderMgtService.GetSender(ctx, senderID)
	if svcErr != nil {
		return svcErr
	}

	if sender.Type != common.NotificationSenderTypeMessage {
		return &ErrorRequestedSenderIsNotOfExpectedType
	}

	_client, svcErr := s.clientFactory.GetClient(ctx, *sender)
	if svcErr != nil {
		return svcErr
	}

	if !_client.IsChannelSupported(channel) {
		return &ErrorUnsupportedChannel
	}

	if err := _client.Send(ctx, channel, data); err != nil {
		s.logger.Error(ctx, "Failed to send notification",
			log.String("channel", string(channel)), log.Error(err))
		return &tidcommon.InternalServerError
	}

	return nil
}
