package service

import (
	"github.com/friendsofgo/errors"

	"tsv-service/internal/services/imports"
	"tsv-service/internal/services/reports"
	"tsv-service/internal/services/units"
)

type ProviderConfiguration struct {
	unitsSvc   *units.Service
	reportsSvc *reports.Service
	importsSvc *imports.Service
}

func NewProviderConfiguration() *ProviderConfiguration {
	return &ProviderConfiguration{}
}

func (p *ProviderConfiguration) SetUnitsService(s *units.Service)     { p.unitsSvc = s }
func (p *ProviderConfiguration) SetReportsService(s *reports.Service) { p.reportsSvc = s }
func (p *ProviderConfiguration) SetImportsService(s *imports.Service) { p.importsSvc = s }

func (p *ProviderConfiguration) UnitsService() (*units.Service, error) {
	if p.unitsSvc == nil {
		return nil, errors.New("units service not set")
	}
	return p.unitsSvc, nil
}

func (p *ProviderConfiguration) ReportsService() (*reports.Service, error) {
	if p.reportsSvc == nil {
		return nil, errors.New("reports service not set")
	}
	return p.reportsSvc, nil
}

func (p *ProviderConfiguration) ImportsService() (*imports.Service, error) {
	if p.importsSvc == nil {
		return nil, errors.New("imports service not set")
	}
	return p.importsSvc, nil
}
