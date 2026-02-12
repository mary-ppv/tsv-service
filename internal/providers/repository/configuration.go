package repository

import (
	"database/sql"
	"tsv-service/internal/repository"

	"github.com/friendsofgo/errors"
)

type ProviderConfiguration struct {
	db *sql.DB

	filesRepo   *repository.FilesRepo
	recordsRepo *repository.RecordsRepo
	reportsRepo *repository.ReportsRepo
	errorsRepo  *repository.ErrorsRepo
	unitsRepo   *repository.UnitsRepo
}

func NewProviderConfiguration(db *sql.DB) *ProviderConfiguration {
	return &ProviderConfiguration{db: db}
}

func (p *ProviderConfiguration) DB() *sql.DB { return p.db }

func (p *ProviderConfiguration) SetFilesRepository(r *repository.FilesRepo)     { p.filesRepo = r }
func (p *ProviderConfiguration) SetRecordsRepository(r *repository.RecordsRepo) { p.recordsRepo = r }
func (p *ProviderConfiguration) SetReportsRepository(r *repository.ReportsRepo) { p.reportsRepo = r }
func (p *ProviderConfiguration) SetErrorsRepository(r *repository.ErrorsRepo)   { p.errorsRepo = r }
func (p *ProviderConfiguration) SetUnitsRepository(r *repository.UnitsRepo)     { p.unitsRepo = r }

func (p *ProviderConfiguration) FilesRepository() (*repository.FilesRepo, error) {
	if p.filesRepo == nil {
		return nil, errors.New("files repository not set")
	}
	return p.filesRepo, nil
}

func (p *ProviderConfiguration) RecordsRepository() (*repository.RecordsRepo, error) {
	if p.recordsRepo == nil {
		return nil, errors.New("records repository not set")
	}
	return p.recordsRepo, nil
}

func (p *ProviderConfiguration) ReportsRepository() (*repository.ReportsRepo, error) {
	if p.reportsRepo == nil {
		return nil, errors.New("reports repository not set")
	}
	return p.reportsRepo, nil
}

func (p *ProviderConfiguration) ErrorsRepository() (*repository.ErrorsRepo, error) {
	if p.errorsRepo == nil {
		return nil, errors.New("errors repository not set")
	}
	return p.errorsRepo, nil
}

func (p *ProviderConfiguration) UnitsRepository() (*repository.UnitsRepo, error) {
	if p.unitsRepo == nil {
		return nil, errors.New("units repository not set")
	}
	return p.unitsRepo, nil
}
