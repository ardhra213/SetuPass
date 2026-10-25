// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authzen

import (
	"encoding/json"
	"net/http"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
)

// handler handles AuthZEN HTTP requests.
type handler struct {
	service AuthZENServiceInterface
}

// errorResponse represents an AuthZEN error response body.
type errorResponse struct {
	Error string `json:"error"`
}

// newHandler creates an AuthZEN HTTP handler.
func newHandler(service AuthZENServiceInterface) *handler {
	return &handler{service: service}
}

// HandleMetadataRequest handles AuthZEN PDP metadata discovery requests.
func (h *handler) HandleMetadataRequest(w http.ResponseWriter, r *http.Request) {
	baseURL := config.GetServerURL(&config.GetServerRuntime().Config.Server)
	resp := MetadataResponse{
		PolicyDecisionPoint:       baseURL,
		AccessEvaluationEndpoint:  baseURL + "/access/v1/evaluation",
		AccessEvaluationsEndpoint: baseURL + "/access/v1/evaluations",
		SearchActionEndpoint:      baseURL + "/access/v1/search/action",
	}

	sysutils.WriteSuccessResponse(r.Context(), w, http.StatusOK, resp)
}

// HandleAccessEvaluationRequest handles a single AuthZEN access evaluation request.
func (h *handler) HandleAccessEvaluationRequest(w http.ResponseWriter, r *http.Request) {
	req, err := sysutils.DecodeJSONBody[AccessEvaluationRequest](r)
	if err != nil {
		handleError(w, &ErrorInvalidRequestFormat)
		return
	}

	resp, svcErr := h.service.EvaluateAccess(r.Context(), *req)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	sysutils.WriteSuccessResponse(r.Context(), w, http.StatusOK, resp)
}

// HandleAccessEvaluationsRequest handles a batched AuthZEN access evaluations request.
func (h *handler) HandleAccessEvaluationsRequest(w http.ResponseWriter, r *http.Request) {
	req, err := sysutils.DecodeJSONBody[AccessEvaluationsRequest](r)
	if err != nil {
		handleError(w, &ErrorInvalidRequestFormat)
		return
	}

	resp, svcErr := h.service.EvaluateAccessBatch(r.Context(), *req)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	sysutils.WriteSuccessResponse(r.Context(), w, http.StatusOK, resp)
}

// HandleActionSearchRequest handles an AuthZEN action search request.
func (h *handler) HandleActionSearchRequest(w http.ResponseWriter, r *http.Request) {
	req, err := sysutils.DecodeJSONBody[AccessActionSearchRequest](r)
	if err != nil {
		handleError(w, &ErrorInvalidRequestFormat)
		return
	}

	resp, svcErr := h.service.SearchActions(r.Context(), *req)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	sysutils.WriteSuccessResponse(r.Context(), w, http.StatusOK, resp)
}

// handleError writes an AuthZEN transport error response for a service error.
func handleError(w http.ResponseWriter, svcErr *tidcommon.ServiceError) {
	statusCode := http.StatusInternalServerError
	if svcErr.Type == tidcommon.ClientErrorType {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Error: svcErr.Error.DefaultValue,
	}); err != nil {
		return
	}
}
