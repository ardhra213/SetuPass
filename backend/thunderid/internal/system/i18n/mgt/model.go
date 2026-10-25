// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package mgt

// --- HTTP Request/Response Models ---

// LanguageListResponse represents the response for listing languages.
type LanguageListResponse struct {
	Languages []string `json:"languages"`
}

// TranslationResponse represents a single translation.
type TranslationResponse struct {
	Language  string `json:"language"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// SetTranslationRequest represents the request body for setting a single translation override.
type SetTranslationRequest struct {
	Value string `json:"value"`
}

// SetTranslationsRequest represents the request body for setting custom translations for a language.
type SetTranslationsRequest struct {
	Translations map[string]map[string]string `json:"translations"`
}

// --- Service Models ---

// Translation represents a translation entity in the service layer.
type Translation struct {
	Key       string `yaml:"id"`
	Language  string `yaml:"language"`
	Namespace string `yaml:"namespace"`
	Value     string `yaml:"value"`
}

// LanguageTranslations represents all translations for a single language.
type LanguageTranslations struct {
	Language     string                       `yaml:"language"`
	Translations map[string]map[string]string `yaml:"translations"`
}
