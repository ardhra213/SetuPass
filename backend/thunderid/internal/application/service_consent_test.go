// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"testing"

	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/application/model"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	"github.com/thunder-id/thunderid/internal/system/log"
)

type ApplicationServiceConsentTestSuite struct {
	suite.Suite
}

func TestApplicationServiceConsentTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationServiceConsentTestSuite))
}

func newTestApplicationService() *applicationService {
	return &applicationService{
		logger: log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ApplicationService")),
	}
}

// ----- validateConsentConfig -----

func (s *ApplicationServiceConsentTestSuite) TestValidateConsentConfig_NilLoginConsent_SetsDefaults() {
	svc := newTestApplicationService()

	app := &model.ApplicationDTO{}
	svc.validateConsentConfig(app)

	s.NotNil(app.LoginConsent)
	s.Equal(int64(0), app.LoginConsent.ValidityPeriod)
}

func (s *ApplicationServiceConsentTestSuite) TestValidateConsentConfig_NegativeValidityPeriodResetToZero() {
	svc := newTestApplicationService()

	app := &model.ApplicationDTO{
		InboundAuthProfile: providers.InboundAuthProfile{
			LoginConsent: &inboundmodel.LoginConsentConfig{ValidityPeriod: -100},
		},
	}
	svc.validateConsentConfig(app)

	s.Equal(int64(0), app.LoginConsent.ValidityPeriod)
}
