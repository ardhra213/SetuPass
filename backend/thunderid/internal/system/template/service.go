// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"context"
	"errors"
	"html"
	"regexp"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"

	"github.com/thunder-id/thunderid/internal/system/log"
)

var ctxPlaceholderRegex = regexp.MustCompile(`\{\{ctx\((\w+)\)}}`)

// templateService implements TemplateServiceInterface using a templateStoreInterface.
type templateService struct {
	store  templateStoreInterface
	logger *log.Logger
}

// newTemplateService creates a new template service with the provided store.
func newTemplateService(store templateStoreInterface) TemplateServiceInterface {
	return &templateService{
		store:  store,
		logger: log.GetLogger().With(log.String(log.LoggerKeyComponentName, "TemplateService")),
	}
}

// GetTemplateByScenario retrieves a template for the specified scenario and template type.
func (s *templateService) GetTemplateByScenario(
	ctx context.Context,
	scenario ScenarioType,
	tmplType TemplateType,
) (*TemplateDTO, *tidcommon.ServiceError) {
	s.logger.Debug(ctx, "Retrieving template by scenario and type",
		log.String("scenario", string(scenario)),
		log.String("type", string(tmplType)))
	tmpl, err := s.store.GetTemplateByScenario(ctx, scenario, tmplType)
	if err != nil {
		if errors.Is(err, errTemplateNotFound) {
			return nil, &ErrorTemplateNotFound
		}
		s.logger.Error(ctx, "Failed to retrieve template by scenario",
			log.String("scenario", string(scenario)),
			log.Error(err))
		return nil, &tidcommon.InternalServerError
	}

	return tmpl, nil
}

// Render renders a template for the specified scenario and template type using the provided data.
func (s *templateService) Render(
	ctx context.Context,
	scenario ScenarioType,
	tmplType TemplateType,
	data TemplateData,
) (*RenderedTemplate, *tidcommon.ServiceError) {
	s.logger.Debug(ctx, "Rendering template", log.String("scenario", string(scenario)))
	tmpl, svcErr := s.GetTemplateByScenario(ctx, scenario, tmplType)
	if svcErr != nil {
		return nil, svcErr
	}

	isHTML := tmpl.ContentType == "text/html"

	replacePlaceholders := func(s string, escapeHTML bool) string {
		return ctxPlaceholderRegex.ReplaceAllStringFunc(s, func(match string) string {
			// Extract the key from {{ctx(key)}}
			submatches := ctxPlaceholderRegex.FindStringSubmatch(match)
			if len(submatches) < 2 {
				return match
			}
			key := submatches[1]
			if val, ok := data[key]; ok {
				if escapeHTML {
					return html.EscapeString(val)
				}
				return val
			}
			return match
		})
	}

	rendered := &RenderedTemplate{

		Subject: replacePlaceholders(tmpl.Subject, false),
		Body:    replacePlaceholders(tmpl.Body, isHTML),
		IsHTML:  isHTML,
	}

	s.logger.Debug(ctx, "Template rendered successfully",
		log.String("scenario", string(scenario)),
		log.String("templateID", tmpl.ID))

	if tmpl.Type == TemplateTypeSMS && len(rendered.Body) > 160 {
		s.logger.Warn(ctx,
			"Rendered SMS body exceeds 160 characters; message may be split into multiple segments",
			log.Int("length", len(rendered.Body)))
	}

	return rendered, nil
}
