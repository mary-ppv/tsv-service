package repository

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"tsv-service/internal/models"
	"tsv-service/internal/repository/driver"
)

type UnitsRepo struct{}

func NewUnitsRepo() *UnitsRepo { return &UnitsRepo{} }

func (r *UnitsRepo) List(ctx context.Context, limit int) (models.UnitSlice, error) {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return nil, err
	}

	mods := []qm.QueryMod{
		qm.OrderBy("created_at desc"),
	}
	if limit > 0 {
		mods = append(mods, qm.Limit(limit))
	}

	return models.Units(mods...).All(ctx, exec)
}
