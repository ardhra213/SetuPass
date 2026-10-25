// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package application

import "strings"

// AppI18nNamespace returns the i18n namespace for application localized metadata.
func AppI18nNamespace() string {
	return "custom"
}

// AppI18nKey returns the i18n key for an application field.
func AppI18nKey(appID, field string) string {
	return "app." + appID + "." + field
}

// AppI18nRef returns the i18n template reference string for an application field.
// The returned value is stored as the application's display field so the UI can
// resolve it to the correct locale at render time.
func AppI18nRef(appID, field string) string {
	return "{{t(" + AppI18nNamespace() + ":" + AppI18nKey(appID, field) + ")}}"
}

// isI18nRef reports whether s is an i18n template reference (e.g. "{{t(custom:app.123.name)}}").
func isI18nRef(s string) bool {
	return strings.HasPrefix(s, "{{t(") && strings.HasSuffix(s, ")}}")
}
