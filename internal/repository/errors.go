package repository

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"

	"tsv-service/internal/models"
	"tsv-service/internal/repository/driver"
)

type ErrorsRepo struct{}

func NewErrorsRepo() *ErrorsRepo { return &ErrorsRepo{} }

func (r *ErrorsRepo) Add(ctx context.Context, e *models.ParseError) error {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return err
	}
	return e.Insert(ctx, exec, boil.Infer())
}

func (r *ErrorsRepo) AddFileLevel(ctx context.Context, fileName, fileID, errText string) error {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return err
	}

	row := &models.ParseError{
		FileName: fileName,
		Error:    errText,
	}

	if fileID != "" {
		row.FileID = null.StringFrom(fileID)
	}

	return row.Insert(ctx, exec, boil.Infer())
}
