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
)

type FilesRepo struct {
	db *sql.DB
}

func NewFilesRepo(db *sql.DB) *FilesRepo { return &FilesRepo{db: db} }

func (r *FilesRepo) EnsureQueued(ctx context.Context, filePath string) (*models.TSVFile, bool, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false, err
	}
	sum := sha256.Sum256(b)
	sha := hex.EncodeToString(sum[:])

	fileName := filepath.Base(filePath)
	
	existing, err := models.TSVFiles(
		models.TSVFileWhere.FileName.EQ(fileName),
	).One(ctx, r.db)

	if err == nil {
		if existing.FileSha256 != sha {
			existing.FileSha256 = sha
			existing.Status = "queued"
			_, _ = existing.Update(ctx, r.db, boil.Infer())
		}
		return existing, false, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, false, err
	}

	row := &models.TSVFile{
		ID:         id.String(),
		FileName:   fileName,
		FileSha256: sha,
		Status:     "queued",
	}

	if err := row.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, false, err
	}

	return row, true, nil
}

func (r *FilesRepo) SetStatus(ctx context.Context, fileID string, status string, errMsg *string) error {
	row, err := models.FindTSVFile(ctx, r.db, fileID)
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

	_, err = row.Update(ctx, r.db, boil.Infer())
	return err
}

func (r *FilesRepo) NextQueued(ctx context.Context, limit int) ([]*models.TSVFile, error) {
	return models.TSVFiles(
		models.TSVFileWhere.Status.EQ("queued"),
		qm.Limit(limit),
		qm.OrderBy("created_at asc"),
	).All(ctx, r.db)
}
