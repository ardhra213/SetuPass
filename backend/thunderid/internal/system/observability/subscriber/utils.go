// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package subscriber

import (
	"context"
	"fmt"

	sysContext "github.com/thunder-id/thunderid/internal/system/context"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/internal/system/observability/event"
	"github.com/thunder-id/thunderid/internal/system/observability/formatter"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Format constants - duplicated here to avoid import cycle with observability package
const (
	formatJSON = "json"
	formatCSV  = "csv"
)

// writerAdapter defines the interface for writing formatted data.
type writerAdapter interface {
	Write(data []byte) error
}

// convertCategories converts string categories to EventCategory types.
// This is a local helper to avoid importing the observability package (which would create a cycle).
func convertCategories(stringCategories []string) []event.EventCategory {
	categories := make([]event.EventCategory, 0, len(stringCategories))
	for _, cat := range stringCategories {
		categories = append(categories, event.EventCategory(cat))
	}
	return categories
}

// processEvent is a shared helper that formats an event and writes it using the provided adapter.
// This eliminates duplicate code between ConsoleSubscriber and FileSubscriber OnEvent implementations.
func processEvent(
	evt *providers.Event,
	fmtr formatter.FormatterInterface,
	adapter writerAdapter,
	logger *log.Logger,
	outputType string,
) error {
	if evt == nil {
		return fmt.Errorf("event is nil")
	}

	// Subscribers run in detached goroutines after the request context may be
	// cancelled, so derive a logging context from the event's trace ID.
	ctx := sysContext.WithTraceID(context.Background(), evt.TraceID)

	// Format the event
	formattedData, err := fmtr.Format(evt)
	if err != nil {
		logger.Error(ctx, "Failed to format event",
			log.String("eventType", evt.Type),
			log.String("eventID", evt.EventID),
			log.Error(err))
		return fmt.Errorf("failed to format event: %w", err)
	}

	// Write using adapter
	if err := adapter.Write(formattedData); err != nil {
		logger.Error(ctx, fmt.Sprintf("Failed to write event to %s", outputType),
			log.String("eventType", evt.Type),
			log.String("eventID", evt.EventID),
			log.Error(err))
		return fmt.Errorf("failed to write to %s: %w", outputType, err)
	}

	logger.Debug(ctx, fmt.Sprintf("Event processed successfully to %s", outputType),
		log.String("eventType", evt.Type),
		log.String("eventID", evt.EventID),
		log.String("traceID", evt.TraceID))

	return nil
}
