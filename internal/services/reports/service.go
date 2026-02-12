package reports

import (
	"context"

	"tsv-service/internal/models"
	"tsv-service/internal/repository"
)

type Service struct {
	repo *repository.ReportsRepo
}

func NewService(repo *repository.ReportsRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListReports(ctx context.Context, unitGUID string, limit int) (models.ReportSlice, error) {
	return s.repo.ListByUnit(ctx, unitGUID, limit)
}
