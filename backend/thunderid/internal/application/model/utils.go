// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package model

// IsValidApplicationType reports whether t is a recognized application type.
func IsValidApplicationType(t ApplicationType) bool {
	switch t {
	case ApplicationTypeBrowser, ApplicationTypeFullStack, ApplicationTypeMobile,
		ApplicationTypeM2M, ApplicationTypeMCP, ApplicationTypeCustom:
		return true
	default:
		return false
	}
}

// ResolveApplicationType converts a raw stored type value into the canonical ApplicationType,
// defaulting to ApplicationTypeCustom when raw is empty or unrecognized. This covers applications
// created before the type attribute existed, or holding a legacy/corrupted value, so callers treat
// them as unconstrained rather than failing.
func ResolveApplicationType(raw string) ApplicationType {
	if t := ApplicationType(raw); IsValidApplicationType(t) {
		return t
	}
	return ApplicationTypeCustom
}
