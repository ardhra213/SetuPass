// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"context"
	"errors"
	"testing"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TemplateServiceTestSuite struct {
	suite.Suite
	mockStore *templateStoreInterfaceMock
	service   TemplateServiceInterface
}

func TestTemplateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateServiceTestSuite))
}

func (suite *TemplateServiceTestSuite) SetupTest() {
	suite.mockStore = newTemplateStoreInterfaceMock(suite.T())
	suite.service = newTemplateService(suite.mockStore)
}

func (suite *TemplateServiceTestSuite) TestGetTemplateByScenario() {
	dto := &TemplateDTO{ID: "test-1", Scenario: ScenarioUserInvite}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite, TemplateTypeEmail).Return(dto, nil)

	res, err := suite.service.GetTemplateByScenario(context.Background(), ScenarioUserInvite, TemplateTypeEmail)
	suite.Nil(err)
	suite.Equal("test-1", res.ID)
}

func (suite *TemplateServiceTestSuite) TestRender() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioUserInvite,
		Subject:     "Test Invite",
		ContentType: "text/html",
		Body:        "Link: {{ctx(inviteLink)}}",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite, TemplateTypeEmail).Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateTypeEmail,
		TemplateData{"inviteLink": "http://example.com"})
	suite.Nil(err)
	suite.Equal("Test Invite", res.Subject)
	suite.Equal("Link: http://example.com", res.Body)
	suite.True(res.IsHTML)
}

func (suite *TemplateServiceTestSuite) TestRender_NotFound() {
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite, TemplateTypeEmail).
		Return(nil, errTemplateNotFound)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateTypeEmail, TemplateData{})
	suite.NotNil(err)
	suite.Equal(&ErrorTemplateNotFound, err)
	suite.Nil(res)
}

func (suite *TemplateServiceTestSuite) TestRender_StoreError() {
	storeErr := errors.New("store error")
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite, TemplateTypeEmail).
		Return(nil, storeErr)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateTypeEmail, TemplateData{})
	suite.NotNil(err)
	suite.Equal(&tidcommon.InternalServerError, err)
	suite.Nil(res)
}

func (suite *TemplateServiceTestSuite) TestRender_UnknownPlaceholder() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioUserInvite,
		Subject:     "Test",
		ContentType: "text/html",
		Body:        "Unknown: {{ctx(unknownKey)}}",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite, TemplateTypeEmail).Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateTypeEmail, TemplateData{})
	suite.Nil(err)
	suite.Equal("Unknown: {{ctx(unknownKey)}}", res.Body)
}

func (suite *TemplateServiceTestSuite) TestRender_EscapesHTMLInBody() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioCIBANotification,
		Subject:     "Authentication Request from {{ctx(appName)}}",
		ContentType: "text/html",
		Body:        "<p>{{ctx(bindingMessage)}}</p>",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioCIBANotification, TemplateTypeEmail).
		Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioCIBANotification, TemplateTypeEmail,
		TemplateData{
			"appName":        "Acme",
			"bindingMessage": `<img src=x onerror=alert(1)>`,
		})
	suite.Nil(err)
	suite.Equal("<p>&lt;img src=x onerror=alert(1)&gt;</p>", res.Body)
}

func (suite *TemplateServiceTestSuite) TestRender_DoesNotEscapeSubject() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioCIBANotification,
		Subject:     "Request from {{ctx(appName)}}",
		ContentType: "text/html",
		Body:        "<p>body</p>",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioCIBANotification, TemplateTypeEmail).
		Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioCIBANotification, TemplateTypeEmail,
		TemplateData{"appName": "AT&T"})
	suite.Nil(err)
	suite.Equal("Request from AT&T", res.Subject)
}

func (suite *TemplateServiceTestSuite) TestRender_DoesNotEscapePlainTextBody() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioCIBANotification,
		Subject:     "Code",
		Type:        TemplateTypeSMS,
		ContentType: "text/plain",
		Body:        "{{ctx(appName)}}: approve?",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioCIBANotification, TemplateTypeSMS).
		Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioCIBANotification, TemplateTypeSMS,
		TemplateData{"appName": "AT&T"})
	suite.Nil(err)
	suite.Equal("AT&T: approve?", res.Body)
	suite.False(res.IsHTML)
}

func (suite *TemplateServiceTestSuite) TestRender_PreservesURLQuerySeparators() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioMagicLink,
		Subject:     "Sign in",
		ContentType: "text/html",
		Body:        `<a href="{{ctx(magicLink)}}">Sign in</a>`,
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioMagicLink, TemplateTypeEmail).
		Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioMagicLink, TemplateTypeEmail,
		TemplateData{"magicLink": "https://localhost:9090/gate/magic?id=abc&token=xyz"})
	suite.Nil(err)
	suite.Equal(`<a href="https://localhost:9090/gate/magic?id=abc&amp;token=xyz">Sign in</a>`, res.Body)
}

func (suite *TemplateServiceTestSuite) TestRender_SubjectPlaceholderReplaced() {
	dto := &TemplateDTO{
		ID:          "1",
		Scenario:    ScenarioSelfRegistration,
		Subject:     "Complete your registration for {{ctx(appName)}}",
		ContentType: "text/html",
		Body:        "Click here: {{ctx(inviteLink)}}",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioSelfRegistration, TemplateTypeEmail).
		Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioSelfRegistration, TemplateTypeEmail,
		TemplateData{"appName": "My App", "inviteLink": "https://example.com/invite"})
	suite.Nil(err)
	suite.Equal("Complete your registration for My App", res.Subject)
	suite.Equal("Click here: https://example.com/invite", res.Body)
}

func (suite *TemplateServiceTestSuite) TestRender_SMSBodyOver160Chars_SucceedsAndReturnsRendered() {
	longBody := "This is an intentionally long SMS body exceeding 160 characters to verify that rendering " +
		"succeeds and returns the rendered result without failing or truncating the content."
	dto := &TemplateDTO{
		ID:          "sms-1",
		Scenario:    ScenarioOTP,
		Type:        TemplateTypeSMS,
		ContentType: "text/plain",
		Body:        longBody,
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioOTP, TemplateTypeSMS).Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioOTP, TemplateTypeSMS, TemplateData{})
	suite.Nil(err)
	suite.Equal(longBody, res.Body)
}

func (suite *TemplateServiceTestSuite) TestRender_EmailBodyOver160Chars_RendersSuccessfully() {
	longBody := "This is an intentionally long email body exceeding 160 characters to verify that no SMS " +
		"segment warning is triggered for non-SMS template types during rendering of the content."
	dto := &TemplateDTO{
		ID:          "email-1",
		Scenario:    ScenarioUserInvite,
		ContentType: "text/html",
		Body:        longBody,
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioUserInvite, TemplateTypeEmail).Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioUserInvite, TemplateTypeEmail, TemplateData{})
	suite.Nil(err)
	suite.Equal(longBody, res.Body)
}

func (suite *TemplateServiceTestSuite) TestRender_SelfRegistrationScenario() {
	dto := &TemplateDTO{
		ID:          "2",
		Scenario:    ScenarioSelfRegistration,
		Subject:     "You're invited",
		ContentType: "text/plain",
		Body:        "Register at {{ctx(inviteLink)}}",
	}
	suite.mockStore.On("GetTemplateByScenario", mock.Anything, ScenarioSelfRegistration, TemplateTypeEmail).
		Return(dto, nil)

	res, err := suite.service.Render(context.Background(), ScenarioSelfRegistration, TemplateTypeEmail,
		TemplateData{"inviteLink": "https://example.com/invite"})
	suite.Nil(err)
	suite.Equal("You're invited", res.Subject)
	suite.Equal("Register at https://example.com/invite", res.Body)
	suite.False(res.IsHTML)
}
