// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package assert

const (
	// AALUnknown represents unknown or no authentication factors.
	AALUnknown AssuranceLevel = "UNKNOWN"
	// AALLevel1 represents basic single-factor authentication.
	AALLevel1 AssuranceLevel = "AAL1"
	// AALLevel2 represents two-factor authentication.
	AALLevel2 AssuranceLevel = "AAL2"
	// AALLevel3 represents multi-factor authentication with hardware token.
	AALLevel3 AssuranceLevel = "AAL3"

	// IALUnknown represents unknown or unverified identity.
	IALUnknown AssuranceLevel = "UNKNOWN"
	// IALLevel1 represents self-asserted identity.
	IALLevel1 AssuranceLevel = "IAL1"
	// IALLevel2 represents identity verified by a trusted party.
	IALLevel2 AssuranceLevel = "IAL2"
	// IALLevel3 represents in-person identity proofing.
	IALLevel3 AssuranceLevel = "IAL3"
)
