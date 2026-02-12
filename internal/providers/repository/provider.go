package repository

import (
	"database/sql"
	"tsv-service/internal/repository"
)

type Provider interface {
	DB() *sql.DB

	FilesRepository() (*repository.FilesRepo, error)
	RecordsRepository() (*repository.RecordsRepo, error)
	ReportsRepository() (*repository.ReportsRepo, error)
	ErrorsRepository() (*repository.ErrorsRepo, error)
	UnitsRepository() (*repository.UnitsRepo, error)
}
