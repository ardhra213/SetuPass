// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"

	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/revocation"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// userDeleteExecutor performs the physical deletion of the subject named in the trusted revocation
// plan. It runs last, after that subject's grants and sessions have been dropped.
type userDeleteExecutor struct {
	providers.Executor
	userProvider userDeletionProvider
}

var _ providers.Executor = (*userDeleteExecutor)(nil)

// newUserDeleteExecutor creates an executor that performs the physical user deletion.
func newUserDeleteExecutor(factory core.FlowFactoryInterface,
	userProvider userDeletionProvider) *userDeleteExecutor {
	base := factory.CreateExecutor(ExecutorNameUserDelete, providers.ExecutorTypeUtility,
		nil, nil, &providers.ExecutorMeta{
			SupportedFlowTypes: []providers.FlowType{providers.FlowTypeAdministration},
		})
	return &userDeleteExecutor{Executor: base, userProvider: userProvider}
}

// Execute deletes the subject identified by the trusted revocation plan.
func (e *userDeleteExecutor) Execute(ctx *providers.NodeContext) (*providers.ExecutorResponse, error) {
	if e.userProvider == nil {
		return nil, fmt.Errorf("user service is not configured")
	}
	plan, err := decodeRevocationPlan(ctx.SharedRuntimeData)
	if err != nil {
		return nil, err
	}
	var subject string
	for _, criterion := range plan.Criteria {
		if criterion.Type == revocation.CriterionTypeSubject {
			subject = criterion.Value
			break
		}
	}
	if subject == "" {
		return nil, fmt.Errorf("trusted revocation plan has no subject")
	}
	if svcErr := e.userProvider.DeleteUser(ctx.Context, subject); svcErr != nil {
		if svcErr.Type == tidcommon.ServerErrorType {
			return nil, fmt.Errorf("failed to delete user: %s", svcErr.Error.DefaultValue)
		}
		return &providers.ExecutorResponse{
			Status: providers.ExecFailure,
			Error:  &ErrUserDeletionFailed,
		}, nil
	}
	return &providers.ExecutorResponse{Status: providers.ExecComplete}, nil
}
