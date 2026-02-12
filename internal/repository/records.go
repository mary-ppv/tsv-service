package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"tsv-service/internal/models"
)

type RecordsRepo struct{ db *sql.DB }

func NewRecordsRepo(db *sql.DB) *RecordsRepo { return &RecordsRepo{db: db} }

func (r *RecordsRepo) Insert(ctx context.Context, rec *models.TSVRecord) error {
	return rec.Insert(ctx, r.db, boil.Infer())
}

func (r *RecordsRepo) ListByUnit(ctx context.Context, unitGuid string, page, limit int) (models.TSVRecordSlice, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	offset := (page - 1) * limit

	total, err := models.TSVRecords(
		models.TSVRecordWhere.UnitGUID.EQ(unitGuid),
	).Count(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	rows, err := models.TSVRecords(
		models.TSVRecordWhere.UnitGUID.EQ(unitGuid),
		qm.OrderBy("created_at desc"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)

	return rows, total, err
}

func (r *RecordsRepo) DeleteByFile(ctx context.Context, fileID string) error {
	_, err := models.TSVRecords(
		models.TSVRecordWhere.FileID.EQ(fileID),
	).DeleteAll(ctx, r.db)
	return err
}
