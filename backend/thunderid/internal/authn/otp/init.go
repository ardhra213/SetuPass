// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package otp

import (
	"github.com/thunder-id/thunderid/internal/notification"
)

// Initialize initializes the OTP authentication service.
func Initialize(notifOTPSvc notification.OTPServiceInterface) OTPAuthnServiceInterface {
	return newOTPAuthnService(notifOTPSvc)
}
