// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package pki

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"path"

	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/http"
	"github.com/thunder-id/thunderid/internal/system/log"
)

// LoadTLSConfig loads a tls.Config from the given certificate and key file paths.
func LoadTLSConfig(cfg *config.Config, certFilePath string, keyFilePath string) (*tls.Config, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "PKIService"))
	// TLS config is loaded at server startup, outside any request,
	// so there is no request context (or trace ID) to propagate.
	ctx := context.Background()

	if certFilePath == "" {
		return nil, errors.New("certificate file path is empty")
	}
	if keyFilePath == "" {
		return nil, errors.New("key file path is empty")
	}

	certFilePath = path.Clean(certFilePath)
	keyFilePath = path.Clean(keyFilePath)

	if _, err := os.Stat(certFilePath); os.IsNotExist(err) {
		return nil, errors.New("certificate file not found at " + certFilePath)
	}
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		return nil, errors.New("key file not found at " + keyFilePath)
	}

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		logger.Error(ctx, "Failed to load X509 key pair", log.Error(err))
		return nil, err
	}

	logger.Debug(ctx, "Successfully loaded TLS certificate",
		log.String("certFile", certFilePath),
		log.String("keyFile", keyFilePath))

	// #nosec G402 -- Min TLS version is TLS 1.2 or higher based on config
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   http.GetTLSVersion(*cfg),
	}, nil
}
