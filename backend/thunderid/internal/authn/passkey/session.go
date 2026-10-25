// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package passkey

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"github.com/thunder-id/thunderid/internal/system/log"
)

const (
	// sessionKeyLength is the length of the random session key in bytes.
	sessionKeyLength = 32
	// sessionTTLSeconds is the session time-to-live in seconds.
	sessionTTLSeconds = 120
)

// generateSessionKey generates a random base64-encoded session key.
func generateSessionKey() (string, error) {
	bytes := make([]byte, sessionKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// storeSessionData stores session data in the database and returns a session key.
func (w *passkeyService) storeSessionData(ctx context.Context,
	sessionData *sessionData,
) (string, *tidcommon.ServiceError) {
	// Generate a random session key
	sessionKey, err := generateSessionKey()
	if err != nil {
		return "", &tidcommon.InternalServerError
	}

	// Store session data in database
	err = w.sessionStore.storeSession(ctx, sessionKey, sessionData, sessionTTLSeconds)
	if err != nil {
		return "", &tidcommon.InternalServerError
	}

	return sessionKey, nil
}

// retrieveSessionData retrieves the session data from the database using the session key.
func (w *passkeyService) retrieveSessionData(ctx context.Context,
	sessionKey string,
) (*sessionData, string, string, *tidcommon.ServiceError) {
	// Retrieve session data from database
	session, err := w.sessionStore.retrieveSession(ctx, sessionKey)
	if err != nil {
		w.logger.Debug(ctx, "Failed to retrieve passkey session", log.Error(err))
		return nil, "", "", &tidcommon.InternalServerError
	}

	if session == nil {
		return nil, "", "", &ErrorSessionExpired
	}

	return session, string(session.UserID), session.RelyingPartyID, nil
}

// clearSessionData removes the session data from the database.
func (w *passkeyService) clearSessionData(ctx context.Context, sessionKey string) {
	// Remove session from database
	_ = w.sessionStore.deleteSession(ctx, sessionKey)
}
