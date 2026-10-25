// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/session"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
	"github.com/thunder-id/thunderid/tests/mocks/flow/sessionmock"
)

type SessionSignOutExecutorTestSuite struct {
	suite.Suite
}

func TestSessionSignOutExecutorTestSuite(t *testing.T) {
	suite.Run(t, new(SessionSignOutExecutorTestSuite))
}

func (suite *SessionSignOutExecutorTestSuite) SetupTest() {
	suite.Require().NoError(config.InitializeServerRuntime(suite.T().TempDir(), &config.Config{}))
}

func (suite *SessionSignOutExecutorTestSuite) TearDownTest() {
	config.ResetServerRuntime()
}

func (suite *SessionSignOutExecutorTestSuite) newExecutor(sso session.Service) *sessionSignOutExecutor {
	flowFactory, _ := core.Initialize(cache.Initialize(config.GetServerRuntime().Config.Cache, "test-deployment"))
	return newSessionSignOutExecutor(flowFactory, sso)
}

// signOutNodeContext carries the login flow's inbound handle and flow id, as the engine delivers them
// through the SSO inputs for a sign-out flow.
func signOutNodeContext() *providers.NodeContext {
	return &providers.NodeContext{
		Context: session.WithSSOInputs(context.Background(), session.SSOInputs{
			Handle: "handle-abc",
			FlowID: "flow-1",
		}),
		ExecutionID: "exec-1",
	}
}

// TestTerminatesAndSignalsClear covers a live session: it is ended and the cookie-clear signal is
// raised on the engine-only channel.
func (suite *SessionSignOutExecutorTestSuite) TestTerminatesAndSignalsClear() {
	sso := sessionmock.NewServiceMock(suite.T())
	sso.EXPECT().Terminate(mock.Anything, "handle-abc", "flow-1").
		Return(&session.Session{SessionID: "sess-1", State: session.StateEnded}, nil)
	exec := suite.newExecutor(sso)

	resp, err := exec.Execute(signOutNodeContext())

	suite.Require().NoError(err)
	suite.Equal(providers.ExecComplete, resp.Status)
	suite.Equal(dataValueTrue, resp.EngineData[common.RuntimeKeySSOSessionCleared])
}

// TestClearsWhenNoSession covers sign-out when no session backs the handle: Terminate is a no-op but
// the cookie is still cleared so the browser drops any stale handle.
func (suite *SessionSignOutExecutorTestSuite) TestClearsWhenNoSession() {
	sso := sessionmock.NewServiceMock(suite.T())
	sso.EXPECT().Terminate(mock.Anything, "handle-abc", "flow-1").Return(nil, nil)
	exec := suite.newExecutor(sso)

	resp, err := exec.Execute(signOutNodeContext())

	suite.Require().NoError(err)
	suite.Equal(providers.ExecComplete, resp.Status)
	suite.Equal(dataValueTrue, resp.EngineData[common.RuntimeKeySSOSessionCleared])
}

// TestTerminateError covers a store failure during termination: the executor surfaces the error and
// does not raise the clear signal.
func (suite *SessionSignOutExecutorTestSuite) TestTerminateError() {
	sso := sessionmock.NewServiceMock(suite.T())
	sso.EXPECT().Terminate(mock.Anything, "handle-abc", "flow-1").Return(nil, errors.New("store down"))
	exec := suite.newExecutor(sso)

	resp, err := exec.Execute(signOutNodeContext())

	suite.Require().Error(err)
	suite.Empty(resp.EngineData[common.RuntimeKeySSOSessionCleared])
}

// TestPromptsWhenConfirmationRequired covers a prompt-enabled node whose logout arrived without a
// valid id_token_hint and no confirmation yet: the executor routes to the confirmation prompt
// (incomplete) and does not terminate the session.
func (suite *SessionSignOutExecutorTestSuite) TestPromptsWhenConfirmationRequired() {
	sso := sessionmock.NewServiceMock(suite.T())
	exec := suite.newExecutor(sso)

	ctx := signOutNodeContext()
	ctx.NodeProperties = map[string]interface{}{propertyKeyPromptOnSignOut: true}
	ctx.RuntimeData = map[string]string{common.RuntimeKeyLogoutPromptRequired: dataValueTrue}

	resp, err := exec.Execute(ctx)

	suite.Require().NoError(err)
	suite.Equal(providers.ExecUserInputRequired, resp.Status)
	suite.Empty(resp.EngineData[common.RuntimeKeySSOSessionCleared])
	sso.AssertNotCalled(suite.T(), "Terminate", mock.Anything, mock.Anything, mock.Anything)
}

// TestTerminatesAfterConfirmation covers the re-run after the End-User confirms: the confirmation
// prompt forwards its confirm action type, so the executor terminates the session instead of
// prompting again.
func (suite *SessionSignOutExecutorTestSuite) TestTerminatesAfterConfirmation() {
	sso := sessionmock.NewServiceMock(suite.T())
	sso.EXPECT().Terminate(mock.Anything, "handle-abc", "flow-1").
		Return(&session.Session{SessionID: "sess-1", State: session.StateEnded}, nil)
	exec := suite.newExecutor(sso)

	ctx := signOutNodeContext()
	ctx.NodeProperties = map[string]interface{}{propertyKeyPromptOnSignOut: true}
	ctx.RuntimeData = map[string]string{common.RuntimeKeyLogoutPromptRequired: dataValueTrue}
	ctx.ForwardedData = map[string]interface{}{
		common.ForwardedDataKeyActionType: string(common.ActionTypeConfirm),
	}

	resp, err := exec.Execute(ctx)

	suite.Require().NoError(err)
	suite.Equal(providers.ExecComplete, resp.Status)
	suite.Equal(dataValueTrue, resp.EngineData[common.RuntimeKeySSOSessionCleared])
}

// TestPromptsWhenActionTypeUnrecognized covers a confirmation prompt that forwarded an action type
// this executor has no case for: an unrecognized type is not consent to end the session, so the
// executor prompts rather than terminating.
func (suite *SessionSignOutExecutorTestSuite) TestPromptsWhenActionTypeUnrecognized() {
	sso := sessionmock.NewServiceMock(suite.T())
	exec := suite.newExecutor(sso)

	ctx := signOutNodeContext()
	ctx.NodeProperties = map[string]interface{}{propertyKeyPromptOnSignOut: true}
	ctx.RuntimeData = map[string]string{common.RuntimeKeyLogoutPromptRequired: dataValueTrue}
	ctx.ForwardedData = map[string]interface{}{
		common.ForwardedDataKeyActionType: "SOME_OTHER_ACTION",
	}

	resp, err := exec.Execute(ctx)

	suite.Require().NoError(err)
	suite.Equal(providers.ExecUserInputRequired, resp.Status)
	suite.Empty(resp.EngineData[common.RuntimeKeySSOSessionCleared])
	sso.AssertNotCalled(suite.T(), "Terminate", mock.Anything, mock.Anything, mock.Anything)
}

// TestTerminatesWhenHintProvided covers a prompt-enabled node whose logout carried a valid
// id_token_hint (no prompt flag): the executor terminates directly without confirming.
func (suite *SessionSignOutExecutorTestSuite) TestTerminatesWhenHintProvided() {
	sso := sessionmock.NewServiceMock(suite.T())
	sso.EXPECT().Terminate(mock.Anything, "handle-abc", "flow-1").Return(nil, nil)
	exec := suite.newExecutor(sso)

	ctx := signOutNodeContext()
	ctx.NodeProperties = map[string]interface{}{propertyKeyPromptOnSignOut: true}

	resp, err := exec.Execute(ctx)

	suite.Require().NoError(err)
	suite.Equal(providers.ExecComplete, resp.Status)
	suite.Equal(dataValueTrue, resp.EngineData[common.RuntimeKeySSOSessionCleared])
}

// TestTerminatesWhenNodeDoesNotOptIn covers a node without the promptOnSignOut property (e.g. the
// always-prompt default flow, where a separate prompt node precedes this one): even when a prompt was
// requested, the executor terminates rather than emitting a second, unhandled prompt.
func (suite *SessionSignOutExecutorTestSuite) TestTerminatesWhenNodeDoesNotOptIn() {
	sso := sessionmock.NewServiceMock(suite.T())
	sso.EXPECT().Terminate(mock.Anything, "handle-abc", "flow-1").Return(nil, nil)
	exec := suite.newExecutor(sso)

	ctx := signOutNodeContext()
	ctx.RuntimeData = map[string]string{common.RuntimeKeyLogoutPromptRequired: dataValueTrue}

	resp, err := exec.Execute(ctx)

	suite.Require().NoError(err)
	suite.Equal(providers.ExecComplete, resp.Status)
	suite.Equal(dataValueTrue, resp.EngineData[common.RuntimeKeySSOSessionCleared])
}
