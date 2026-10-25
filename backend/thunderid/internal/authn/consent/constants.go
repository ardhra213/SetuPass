// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package consent

const (
	// consentSessionTokenAudience is the JWT audience for consent session tokens
	consentSessionTokenAudience = "consent-svc"

	// consentSessionTokenValidityPeriod is the validity period of consent session tokens in seconds
	consentSessionTokenValidityPeriod = int64(3600)

	// consentSessionClaimKey is the JWT claim key for consent session data
	consentSessionClaimKey = "consent_session"
)
