package repository

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"tsv-service/internal/models"
	"tsv-service/internal/repository/driver"
)

type ReportsRepo struct{}

func NewReportsRepo() *ReportsRepo { return &ReportsRepo{} }

func (r *ReportsRepo) Insert(ctx context.Context, report *models.Report) error {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return err
	}
	return report.Insert(ctx, exec, boil.Infer())
}

func (r *ReportsRepo) ListByUnit(ctx context.Context, unitGUID string, limit int) (models.ReportSlice, error) {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return nil, err
	}

	mods := []qm.QueryMod{
		models.ReportWhere.UnitGUID.EQ(unitGUID),
		qm.OrderBy("created_at desc"),
	}
	if limit > 0 {
		mods = append(mods, qm.Limit(limit))
	}

	return models.Reports(mods...).All(ctx, exec)
}
