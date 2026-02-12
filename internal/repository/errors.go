package repository

import (
	"context"
	"database/sql"

	"tsv-service/internal/models"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
)

type ErrorsRepo struct{ db *sql.DB }

func NewErrorsRepo(db *sql.DB) *ErrorsRepo { return &ErrorsRepo{db: db} }

func (r *ErrorsRepo) Add(ctx context.Context, e *models.ParseError) error {
	return e.Insert(ctx, r.db, boil.Infer())
}

func (r *ErrorsRepo) AddFileLevel(
	ctx context.Context,
	fileName string,
	fileID string,
	errText string,
) error {

	row := &models.ParseError{
		FileName: fileName,
		Error:    errText,
	}

	if fileID != "" {
		row.FileID = null.StringFrom(fileID)
	}

	return row.Insert(ctx, r.db, boil.Infer())
}
