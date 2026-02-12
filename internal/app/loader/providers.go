package loader

import (
	"context"
	"database/sql"
	repoProvider "tsv-service/internal/providers/repository"
	svcProvider "tsv-service/internal/providers/service"

	"tsv-service/internal/repository"
	"tsv-service/internal/services/imports"
	"tsv-service/internal/services/reports"
	"tsv-service/internal/services/units"
)

func InitRepositoryProvider(db *sql.DB) repoProvider.Provider {
	p := repoProvider.NewProviderConfiguration(db)

	p.SetFilesRepository(repository.NewFilesRepo())
	p.SetRecordsRepository(repository.NewRecordsRepo())
	p.SetErrorsRepository(repository.NewErrorsRepo())
	p.SetReportsRepository(repository.NewReportsRepo())
	p.SetUnitsRepository(repository.NewUnitsRepo())

	return p
}

func InitServiceProvider(
	_ context.Context,
	rp repoProvider.Provider,
) svcProvider.Provider {
	p := svcProvider.NewProviderConfiguration()

	unitsRepo, _ := rp.UnitsRepository()
	recordsRepo, _ := rp.RecordsRepository()
	reportsRepo, _ := rp.ReportsRepository()

	p.SetUnitsService(units.NewService(unitsRepo, recordsRepo))
	p.SetReportsService(reports.NewService(reportsRepo))
	p.SetImportsService(imports.NewService(unitsRepo))

	return p
}
