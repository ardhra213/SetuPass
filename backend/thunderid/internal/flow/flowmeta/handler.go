// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package flowmeta

import (
	"context"
	"net/http"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"github.com/thunder-id/thunderid/internal/system/error/apierror"
	"github.com/thunder-id/thunderid/internal/system/log"
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
)

// flowMetaHandler handles flow metadata HTTP requests.
type flowMetaHandler struct {
	flowMetaService FlowMetaServiceInterface
	logger          *log.Logger
}

// newFlowMetaHandler creates a new instance of flowMetaHandler.
func newFlowMetaHandler(flowMetaService FlowMetaServiceInterface) *flowMetaHandler {
	return &flowMetaHandler{
		flowMetaService: flowMetaService,
		logger:          log.GetLogger().With(log.String(log.LoggerKeyComponentName, "FlowMetaHandler")),
	}
}

// HandleGetFlowMetadata handles the GET /flow/meta endpoint.
func (h *flowMetaHandler) HandleGetFlowMetadata(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	metaType := sysutils.SanitizeString(r.URL.Query().Get("type"))
	id := sysutils.SanitizeString(r.URL.Query().Get("id"))

	var language *string
	var namespace *string

	if lang := r.URL.Query().Get("language"); lang != "" {
		language = &lang
	}

	if ns := r.URL.Query().Get("namespace"); ns != "" {
		namespace = &ns
	}

	// Validate parameter combinations: id requires type, and type requires id
	if id != "" && metaType == "" {
		handleServiceError(r.Context(), w, &ErrorMissingType)
		return
	}

	if metaType != "" && id == "" {
		handleServiceError(r.Context(), w, &ErrorMissingID)
		return
	}
	if language != nil {
		lang := sysutils.SanitizeString(*language)
		language = &lang
	}
	if namespace != nil {
		ns := sysutils.SanitizeString(*namespace)
		namespace = &ns
	}

	// Call service
	metadata, svcErr := h.flowMetaService.GetFlowMetadata(r.Context(), MetaType(metaType), id, language, namespace)
	if svcErr != nil {
		handleServiceError(r.Context(), w, svcErr)
		return
	}

	// Return success response
	sysutils.WriteSuccessResponse(r.Context(), w, http.StatusOK, metadata)
	h.logger.Debug(r.Context(), "Flow metadata retrieved successfully",
		log.String("type", metaType),
		log.String("id", id))
}

// handleServiceError converts service errors to appropriate HTTP responses.
func handleServiceError(ctx context.Context, w http.ResponseWriter, svcErr *tidcommon.ServiceError) {
	errResp := apierror.ErrorResponse{
		Code:        svcErr.Code,
		Message:     svcErr.Error,
		Description: svcErr.ErrorDescription,
	}

	statusCode := http.StatusInternalServerError
	if svcErr.Type == tidcommon.ClientErrorType {
		// Determine specific client error status code
		if svcErr.Code == ErrorApplicationNotFound.Code || svcErr.Code == ErrorOUNotFound.Code {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusBadRequest
		}
	}

	sysutils.WriteErrorResponse(ctx, w, statusCode, errResp)
}
