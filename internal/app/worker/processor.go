package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/gofrs/uuid/v5"

	"tsv-service/internal/models"
	"tsv-service/internal/repository"
)

type Processor struct {
	db          *sql.DB
	filesRepo   *repository.FilesRepo
	recordsRepo *repository.RecordsRepo
	errorsRepo  *repository.ErrorsRepo
	reportGen   *ReportGenerator

	inputDir string
}

func NewProcessor(
	db *sql.DB,
	inputDir string,
	filesRepo *repository.FilesRepo,
	recordsRepo *repository.RecordsRepo,
	errorsRepo *repository.ErrorsRepo,
	reportGen *ReportGenerator,
) *Processor {
	return &Processor{
		db:          db,
		inputDir:    inputDir,
		filesRepo:   filesRepo,
		recordsRepo: recordsRepo,
		errorsRepo:  errorsRepo,
		reportGen:   reportGen,
	}
}

func (p *Processor) Process(ctx context.Context, file *models.TSVFile) (retErr error) {
	fullPath := filepath.Join(p.inputDir, file.FileName)

	_ = p.filesRepo.SetStatus(ctx, file.ID, "processing", nil)

	defer func() {
		if retErr != nil {
			msg := retErr.Error()
			_ = p.filesRepo.SetStatus(ctx, file.ID, "failed", &msg)
			_ = p.errorsRepo.AddFileLevel(ctx, file.FileName, file.ID, msg)
			return
		}
		_ = p.filesRepo.SetStatus(ctx, file.ID, "done", nil)
	}()

	lines, parseErrs, err := ParseTSVFile(fullPath)
	if err != nil {
		msg := err.Error()
		_ = p.filesRepo.SetStatus(ctx, file.ID, "failed", &msg)
		_ = p.errorsRepo.AddFileLevel(ctx, file.FileName, file.ID, msg)
		return err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	exec := tx

	unitToRows := map[string][]models.TSVRecord{}

	for _, pe := range parseErrs {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		row := &models.ParseError{
			ID:       id.String(),
			FileID:   null.StringFrom(file.ID),
			FileName: file.FileName,
			Error:    pe.Err.Error(),
		}

		if pe.RawLine != "" {
			row.RawLine = null.StringFrom(pe.RawLine)
		}
		if pe.LineNo > 0 {
			row.LineNo = null.IntFrom(pe.LineNo)
		}
		if pe.UnitGuid != nil && *pe.UnitGuid != "" {
			row.UnitGUID = null.StringFrom(*pe.UnitGuid)
		}

		if err := row.Insert(ctx, exec, boil.Infer()); err != nil {
			return err
		}
	}

	for i := range lines {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		rec := &models.TSVRecord{
			ID:       id.String(),
			FileID:   file.ID,
			UnitGUID: lines[i].UnitGuid,
			LineNo:   lines[i].LineNo,
			Payload:  types.JSON(mustJSON(lines[i].Payload)),
		}

		if strings.TrimSpace(lines[i].MsgID) != "" {
			rec.MSGID = null.StringFrom(lines[i].MsgID)
		}

		if lines[i].MsgID != "" {
			rec.MSGID = null.StringFrom(lines[i].MsgID)
		}

		if err := rec.Insert(ctx, exec, boil.Infer()); err != nil {
			return err
		}

		unitToRows[rec.UnitGUID] = append(unitToRows[rec.UnitGUID], *rec)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for unit, rows := range unitToRows {
		if _, err := p.reportGen.WriteUnitReport(ctx, unit, file.ID, rows); err != nil {
			_ = p.errorsRepo.AddFileLevel(ctx, file.FileName, file.ID,
				fmt.Sprintf("report failed for %s: %v", unit, err),
			)
		}
	}

	_ = p.filesRepo.SetStatus(ctx, file.ID, "done", nil)
	return nil
}

func mustJSON(m map[string]any) []byte {
	b, _ := json.Marshal(m)
	return b
}
