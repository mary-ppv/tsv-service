package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"tsv-service/internal/models"
)

type ReportsRepo struct {
	db *sql.DB
}

func NewReportsRepo(db *sql.DB) *ReportsRepo {
	return &ReportsRepo{db: db}
}

func (r *ReportsRepo) Insert(ctx context.Context, report *models.Report) error {
	return report.Insert(ctx, r.db, boil.Infer())
}

func (r *ReportsRepo) ListByUnit(ctx context.Context, unitGUID string, limit int) (models.ReportSlice, error) {
	mods := []qm.QueryMod{
		models.ReportWhere.UnitGUID.EQ(unitGUID),
		qm.OrderBy("created_at desc"),
	}
	if limit > 0 {
		mods = append(mods, qm.Limit(limit))
	}

	return models.Reports(mods...).All(ctx, r.db)
}
