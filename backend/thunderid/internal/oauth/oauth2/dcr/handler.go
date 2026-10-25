// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package dcr

import (
	"context"
	"net/http"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	oauthconfig "github.com/thunder-id/thunderid/internal/oauth/config"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/internal/system/security"
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
)

// dcrHandler defines the handler for DCR API requests.
type dcrHandler struct {
	dcrInsecure bool
	dcrService  DCRServiceInterface
}

// newDCRHandler creates a new instance of dcrHandler.
func newDCRHandler(dcrService DCRServiceInterface, cfg oauthconfig.Config) *dcrHandler {
	return &dcrHandler{
		dcrInsecure: cfg.OAuth.DCR.Insecure,
		dcrService:  dcrService,
	}
}

// HandleDCRRegistration handles the DCR client registration request.
func (dh *dcrHandler) HandleDCRRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// When DCR is not insecure, require a valid token with required permissions.
	if !dh.dcrInsecure && !dh.checkDCRAuthorization(r, w) {
		return
	}

	dcrRequest, err := sysutils.DecodeJSONBody[DCRRegistrationRequest](r)
	if err != nil {
		sysutils.WriteJSONError(ctx, w, ErrorInvalidRequestFormat.Code,
			ErrorInvalidRequestFormat.ErrorDescription.DefaultValue, http.StatusBadRequest, nil)
		return
	}

	dcrResponse, svcErr := dh.dcrService.RegisterClient(ctx, dcrRequest)
	if svcErr != nil {
		if svcErr.Type == tidcommon.ServerErrorType {
			logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "DCRHandler"))
			logger.Error(ctx, "Internal server error processing DCR registration request",
				log.MaskedString("client_name", dcrRequest.ClientName),
				log.String("error_code", svcErr.Code),
				log.String("error", svcErr.Error.DefaultValue),
			)
		}
		dh.writeServiceErrorResponse(ctx, w, svcErr)
		return
	}

	sysutils.WriteSuccessResponse(ctx, w, http.StatusCreated, dcrResponse)
}

// checkDCRAuthorization verifies that the caller holds required permission.
// Returns true if authorized, false (and writes an HTTP 401) otherwise.
func (dh *dcrHandler) checkDCRAuthorization(r *http.Request, w http.ResponseWriter) bool {
	if security.HasSystemPermission(security.GetPermissions(r.Context())) {
		return true
	}
	sysutils.WriteJSONError(r.Context(), w, ErrorUnauthorized.Code,
		ErrorUnauthorized.ErrorDescription.DefaultValue, http.StatusUnauthorized, nil)
	return false
}

// writeServiceErrorResponse writes a service error response.
func (
	dh *dcrHandler) writeServiceErrorResponse(ctx context.Context,
	w http.ResponseWriter,
	svcErr *tidcommon.ServiceError) {
	var statusCode int

	switch svcErr.Type {
	case tidcommon.ClientErrorType:
		statusCode = http.StatusBadRequest
	case tidcommon.ServerErrorType:
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusBadRequest
	}

	sysutils.WriteJSONError(ctx, w, svcErr.Code, svcErr.ErrorDescription.DefaultValue, statusCode, nil)
}
