// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package template

// TemplateType represents the type of template (e.g., email, sms).
type TemplateType string

const (
	// TemplateTypeEmail represents an email template.
	TemplateTypeEmail TemplateType = "email"
	// TemplateTypeSMS represents an SMS template.
	TemplateTypeSMS TemplateType = "sms"
)

// ScenarioType represents the scenario for which a template is used.
type ScenarioType string

const (
	// ScenarioUserInvite represents the user invitation scenario.
	ScenarioUserInvite ScenarioType = "USER_INVITE"
	// ScenarioMagicLink represents the magic link sign-in scenario.
	ScenarioMagicLink ScenarioType = "MAGIC_LINK"
	// ScenarioSelfRegistration represents the self-registration via invite link scenario.
	ScenarioSelfRegistration ScenarioType = "SELF_REGISTRATION"
	// ScenarioOTP represents the OTP verification scenario.
	ScenarioOTP ScenarioType = "OTP"
	// ScenarioPasswordRecovery represents the password recovery via email link scenario.
	ScenarioPasswordRecovery ScenarioType = "PASSWORD_RECOVERY"
	// ScenarioCIBANotification represents the CIBA backchannel authentication notification scenario.
	ScenarioCIBANotification ScenarioType = "CIBA_NOTIFICATION"
)

// supportedScenarios contains all valid scenario types.
var supportedScenarios = map[ScenarioType]bool{
	ScenarioUserInvite:       true,
	ScenarioMagicLink:        true,
	ScenarioSelfRegistration: true,
	ScenarioOTP:              true,
	ScenarioPasswordRecovery: true,
	ScenarioCIBANotification: true,
}

// IsValidScenario checks if the given scenario type is supported.
func IsValidScenario(scenario ScenarioType) bool {
	return supportedScenarios[scenario]
}

// TemplateDTO represents a template with embedded metadata.
type TemplateDTO struct {
	ID          string       `yaml:"id"`
	DisplayName string       `yaml:"displayName"`
	Scenario    ScenarioType `yaml:"scenario"`
	Type        TemplateType `yaml:"type"`
	Subject     string       `yaml:"subject"`
	ContentType string       `yaml:"contentType"`
	Body        string       `yaml:"body"`
}

// TemplateData holds key-value pairs for template substitution.
type TemplateData = map[string]string

// RenderedTemplate holds the result after template processing.
type RenderedTemplate struct {
	Subject string
	Body    string
	IsHTML  bool
}
