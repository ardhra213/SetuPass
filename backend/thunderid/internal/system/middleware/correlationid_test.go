// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sysContext "github.com/thunder-id/thunderid/internal/system/context"
)

func TestCorrelationIDMiddleware_GeneratesID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := sysContext.GetTraceID(r.Context())
		if traceID == "" {
			t.Error("Expected trace ID in context, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// Check response header
	correlationID := w.Header().Get("X-Correlation-ID")
	if correlationID == "" {
		t.Error("Expected X-Correlation-ID header in response, got empty")
	}
}

func TestCorrelationIDMiddleware_ExtractsFromXCorrelationID(t *testing.T) {
	expectedID := "test-correlation-id-123"
	var actualID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualID = sysContext.GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Correlation-ID", expectedID)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if actualID != expectedID {
		t.Errorf("Expected trace ID %s, got %s", expectedID, actualID)
	}

	// Check response header matches
	responseID := w.Header().Get("X-Correlation-ID")
	if responseID != expectedID {
		t.Errorf("Expected response header %s, got %s", expectedID, responseID)
	}
}

func TestCorrelationIDMiddleware_ExtractsFromXRequestID(t *testing.T) {
	expectedID := "test-request-id-456"
	var actualID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualID = sysContext.GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", expectedID)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if actualID != expectedID {
		t.Errorf("Expected trace ID %s, got %s", expectedID, actualID)
	}
}

func TestCorrelationIDMiddleware_ExtractsFromXTraceID(t *testing.T) {
	expectedID := "test-trace-id-789"
	var actualID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualID = sysContext.GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Trace-ID", expectedID)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if actualID != expectedID {
		t.Errorf("Expected trace ID %s, got %s", expectedID, actualID)
	}
}

func TestCorrelationIDMiddleware_PriorityOrder(t *testing.T) {
	// X-Correlation-ID should take priority over others
	expectedID := "correlation-id"
	var actualID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualID = sysContext.GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Correlation-ID", expectedID)
	req.Header.Set("X-Request-ID", "request-id")
	req.Header.Set("X-Trace-ID", "trace-id")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if actualID != expectedID {
		t.Errorf("Expected trace ID %s (X-Correlation-ID should take priority), got %s", expectedID, actualID)
	}
}

func TestCorrelationIDMiddleware_XRequestIDPriorityOverXTraceID(t *testing.T) {
	// X-Request-ID should take priority over X-Trace-ID
	expectedID := "request-id"
	var actualID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualID = sysContext.GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", expectedID)
	req.Header.Set("X-Trace-ID", "trace-id")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if actualID != expectedID {
		t.Errorf("Expected trace ID %s (X-Request-ID should take priority over X-Trace-ID), got %s",
			expectedID, actualID)
	}
}

func TestCorrelationIDMiddleware_ResponseHeaderAlwaysSet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorrelationIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	correlationID := w.Header().Get("X-Correlation-ID")
	if correlationID == "" {
		t.Error("Expected X-Correlation-ID header to be set in response")
	}
}

func TestExtractCorrelationID_NoHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	id := extractCorrelationID(req)

	if id != "" {
		t.Errorf("Expected empty string when no headers present, got %s", id)
	}
}

func TestExtractCorrelationID_WithXCorrelationID(t *testing.T) {
	expectedID := "test-id"
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Correlation-ID", expectedID)

	id := extractCorrelationID(req)
	if id != expectedID {
		t.Errorf("Expected %s, got %s", expectedID, id)
	}
}
