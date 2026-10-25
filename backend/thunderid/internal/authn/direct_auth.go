// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package authn

import (
	"crypto/subtle"
	"net/http"

	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	"github.com/thunder-id/thunderid/internal/system/error/apierror"
	sysutils "github.com/thunder-id/thunderid/internal/system/utils"
)

// directAuthHeaderName is the request header carrying the Direct Auth Secret on Direct API requests.
const directAuthHeaderName = "Direct-Auth-Secret"

// DirectAuthGuardInterface enforces the Direct Auth Secret on the Direct API endpoints
// (/auth/**, /register/passkey/**, /access/**).
type DirectAuthGuardInterface interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}

// directAuthGuard is secure by default: an empty secret blocks every wrapped endpoint.
type directAuthGuard struct {
	secret string
}

// newDirectAuthGuard creates the Direct Auth Secret guard.
func newDirectAuthGuard(secret string) DirectAuthGuardInterface {
	return &directAuthGuard{secret: secret}
}

// Wrap admits the request only when the configured secret matches the Direct-Auth-Secret header
// (constant-time compare); otherwise it responds 401 with an RFC 6750 Bearer challenge.
func (g *directAuthGuard) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provided := r.Header.Get(directAuthHeaderName)
		if g.secret == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(g.secret)) != 1 {
			g.writeUnauthorized(w, r)
			return
		}
		next(w, r)
	}
}

func (g *directAuthGuard) writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(serverconst.WWWAuthenticateHeaderName, serverconst.TokenTypeBearer)
	sysutils.WriteErrorResponse(r.Context(), w, http.StatusUnauthorized, apierror.ErrUnauthorized)
}
