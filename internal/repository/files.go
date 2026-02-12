package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/friendsofgo/errors"
	"github.com/gofrs/uuid/v5"

	"tsv-service/internal/models"
	"tsv-service/internal/repository/driver"
)

type FilesRepo struct{}

func NewFilesRepo() *FilesRepo { return &FilesRepo{} }

func (r *FilesRepo) EnsureQueued(ctx context.Context, filePath string) (*models.TSVFile, bool, error) {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return nil, false, err
	}

	fileName := filepath.Base(filePath)

	existing, err := models.TSVFiles(
		models.TSVFileWhere.FileName.EQ(fileName),
	).One(ctx, exec)

	if err == nil {
		return existing, false, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false, err
	}
	sum := sha256.Sum256(b)
	sha := hex.EncodeToString(sum[:])

	id, _ := uuid.NewV7()

	row := &models.TSVFile{
		ID:         id.String(),
		FileName:   fileName,
		FileSha256: sha,
		Status:     "queued",
	}

	if err := row.Insert(ctx, exec, boil.Infer()); err != nil {
		return nil, false, err
	}

	return row, true, nil
}

func (r *FilesRepo) SetStatus(ctx context.Context, fileID string, status string, errMsg *string) error {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return err
	}

	row, err := models.FindTSVFile(ctx, exec, fileID)
	if err != nil {
		return err
	}

	row.Status = status

	if errMsg != nil {
		row.ErrorMessage.Valid = true
		row.ErrorMessage.String = *errMsg
	} else {
		row.ErrorMessage.Valid = false
		row.ErrorMessage.String = ""
	}

	_, err = row.Update(ctx, exec, boil.Infer())
	return err
}

func (r *FilesRepo) NextQueued(ctx context.Context, limit int) (models.TSVFileSlice, error) {
	exec, err := driver.ExecutorFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return models.TSVFiles(
		models.TSVFileWhere.Status.EQ("queued"),
		qm.Limit(limit),
		qm.OrderBy("created_at asc"),
	).All(ctx, exec)
}
