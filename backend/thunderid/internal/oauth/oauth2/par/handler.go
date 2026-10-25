// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package par

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/clientauth"
	oauth2const "github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/dpop"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/internal/system/utils"
)

// parHandlerInterface defines the interface for handling PAR requests.
type parHandlerInterface interface {
	HandlePARRequest(w http.ResponseWriter, r *http.Request)
}

// parHandler implements parHandlerInterface.
type parHandler struct {
	parService   PARServiceInterface
	dpopVerifier dpop.VerifierInterface
	parEndpoint  string
	logger       *log.Logger
}

// newPARHandler creates a new PAR handler instance.
func newPARHandler(
	parService PARServiceInterface, dpopVerifier dpop.VerifierInterface, parEndpoint string,
) parHandlerInterface {
	return &parHandler{
		parService:   parService,
		dpopVerifier: dpopVerifier,
		parEndpoint:  parEndpoint,
		logger:       log.GetLogger().With(log.String(log.LoggerKeyComponentName, "PARHandler")),
	}
}

// HandlePARRequest handles the POST /oauth2/par request.
func (h *parHandler) HandlePARRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Client authentication is handled by the ClientAuthMiddleware.
	clientInfo := clientauth.GetOAuthClient(ctx)
	if clientInfo == nil {
		h.logger.Error(ctx, "OAuth client not found in context - ClientAuthMiddleware must be applied")
		utils.WriteJSONError(ctx, w, oauth2const.ErrorServerError,
			"Something went wrong", http.StatusInternalServerError, nil)
		return
	}

	// Parse form-encoded body.
	if err := r.ParseForm(); err != nil {
		utils.WriteJSONError(ctx, w, oauth2const.ErrorInvalidRequest, "Failed to parse request body",
			http.StatusBadRequest, nil)
		return
	}

	// A DPoP proof at the PAR endpoint binds the auth code to the proof's key.
	dpopHeaderJkt, errCode, errDesc := h.verifyDPoPHeader(ctx, r)
	if errCode != "" {
		statusCode := http.StatusBadRequest
		if errCode == oauth2const.ErrorServerError {
			h.logger.Error(ctx, "Internal server error verifying DPoP header",
				log.MaskedString("clientID", clientInfo.ClientID),
				log.String("errorCode", errCode),
				log.String("errorDescription", errDesc),
			)
			statusCode = http.StatusInternalServerError
		}
		if errCode == oauth2const.ErrorInvalidDPoPProof {
			h.logger.Debug(ctx, "DPoP proof rejected at PAR",
				log.MaskedString("clientID", clientInfo.ClientID),
				log.String("error", errDesc))
			errDesc = "Invalid DPoP proof"
		}
		utils.WriteJSONError(ctx, w, errCode, errDesc, statusCode, nil)
		return
	}

	// Only the resource parameter is permitted to be repeated (RFC 8707 §2). Any other parameter
	// appearing more than once is rejected with invalid_request per RFC 6749 §3.1.
	params := make(map[string]string, len(r.PostForm))
	for key, values := range r.PostForm {
		if key == oauth2const.RequestParamResource {
			continue
		}
		if len(values) > 1 {
			utils.WriteJSONError(ctx, w, oauth2const.ErrorInvalidRequest,
				fmt.Sprintf("Parameter %q must not be repeated", key), http.StatusBadRequest, nil)
			return
		}
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	resources := r.PostForm[oauth2const.RequestParamResource]

	resp, errCode, errDesc := h.parService.HandlePushedAuthorizationRequest(
		ctx, params, resources, clientInfo.OAuthApp, dpopHeaderJkt)
	if errCode != "" {
		statusCode := http.StatusBadRequest
		if errCode == oauth2const.ErrorServerError {
			h.logger.Error(ctx, "Internal server error processing pushed authorization request",
				log.MaskedString("clientID", clientInfo.ClientID),
				log.String("errorCode", errCode),
				log.String("errorDescription", errDesc),
			)
			statusCode = http.StatusInternalServerError
		}
		utils.WriteJSONError(ctx, w, errCode, errDesc, statusCode, nil)
		return
	}

	utils.WriteSuccessResponse(ctx, w, http.StatusCreated, resp)
}

// verifyDPoPHeader verifies the DPoP proof header if present and returns the JKT, error code, and error description.
func (h *parHandler) verifyDPoPHeader(ctx context.Context, r *http.Request) (string, string, string) {
	dpopHeaders := r.Header.Values(oauth2const.HeaderDPoP)
	if len(dpopHeaders) == 0 {
		return "", "", ""
	}
	if len(dpopHeaders) > 1 {
		return "", oauth2const.ErrorInvalidDPoPProof, "Multiple DPoP headers"
	}
	if h.dpopVerifier == nil {
		h.logger.Error(ctx, "DPoP verifier is not configured")
		return "", oauth2const.ErrorServerError, "Something went wrong"
	}
	result, err := h.dpopVerifier.Verify(ctx, dpop.VerifyParams{
		Proof: dpopHeaders[0],
		HTM:   http.MethodPost,
		HTU:   h.parEndpoint,
	})
	if err != nil {
		if errors.Is(err, dpop.ErrReplayedProof) {
			return "", oauth2const.ErrorInvalidDPoPProof, "DPoP proof replayed"
		}
		return "", oauth2const.ErrorInvalidDPoPProof, err.Error()
	}
	return result.JKT, "", ""
}
