package repository

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/gofrs/uuid/v5"

	"tsv-service/internal/models"
	"tsv-service/internal/repository/driver"
)

type RecordsRepo struct{}

func NewRecordsRepo() *RecordsRepo { return &RecordsRepo{} }

func (r *RecordsRepo) Insert(ctx context.Context, rec *models.TSVRecord) error {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return err
	}
	return rec.Insert(ctx, exec, boil.Infer())
}

func (r *RecordsRepo) ListByUnit(ctx context.Context, unitGUID string, page, limit int) (models.TSVRecordSlice, int64, error) {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	offset := (page - 1) * limit

	total, err := models.TSVRecords(
		models.TSVRecordWhere.UnitGUID.EQ(unitGUID),
	).Count(ctx, exec)
	if err != nil {
		return nil, 0, err
	}

	rows, err := models.TSVRecords(
		models.TSVRecordWhere.UnitGUID.EQ(unitGUID),
		qm.OrderBy("created_at desc"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, exec)

	return rows, total, err
}

func (r *RecordsRepo) DeleteByFile(ctx context.Context, fileID uuid.UUID) error {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = models.TSVRecords(
		models.TSVRecordWhere.FileID.EQ(fileID.String()),
	).DeleteAll(ctx, exec)

	return err
}
