package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"tsv-service/internal/models"
)

type UnitsRepo struct{ db *sql.DB }

func NewUnitsRepo(db *sql.DB) *UnitsRepo { return &UnitsRepo{db: db} }

func (r *UnitsRepo) List(ctx context.Context, limit int) (models.UnitSlice, error) {
	mods := []qm.QueryMod{qm.OrderBy("created_at desc")}
	if limit > 0 {
		mods = append(mods, qm.Limit(limit))
	}
	return models.Units(mods...).All(ctx, r.db)
}
