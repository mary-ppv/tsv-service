package imports

import (
	"context"

	"tsv-service/internal/repository"
)

type Service struct {
	unitsRepo *repository.UnitsRepo
}

func NewService(unitsRepo *repository.UnitsRepo) *Service {
	return &Service{unitsRepo: unitsRepo}
}

func (s *Service) ImportSourceData(ctx context.Context, path string) error {
	// TODO: parse file -> upsert units / etc
	return nil
}
