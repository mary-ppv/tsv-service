package units

import (
	"context"

	"tsv-service/internal/models"
	"tsv-service/internal/repository"
)

type Service struct {
	unitsRepo   *repository.UnitsRepo
	recordsRepo *repository.RecordsRepo
}

func NewService(unitsRepo *repository.UnitsRepo, recordsRepo *repository.RecordsRepo) *Service {
	return &Service{unitsRepo: unitsRepo, recordsRepo: recordsRepo}
}

func (s *Service) ListRecords(ctx context.Context, unitGUID string, page, limit int) (models.TSVRecordSlice, int64, error) {
	return s.recordsRepo.ListByUnit(ctx, unitGUID, page, limit)
}
