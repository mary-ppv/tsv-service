package service

import (
	"tsv-service/internal/services/imports"
	"tsv-service/internal/services/reports"
	"tsv-service/internal/services/units"
)

type Provider interface {
	UnitsService() (*units.Service, error)
	ReportsService() (*reports.Service, error)
	ImportsService() (*imports.Service, error)
}
