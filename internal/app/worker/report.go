package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"tsv-service/internal/models"
	"tsv-service/internal/repository"

	"github.com/gofrs/uuid/v5"
)

type ReportGenerator struct {
	outDir      string
	reportsRepo *repository.ReportsRepo
}

func NewReportGenerator(outDir string, repo *repository.ReportsRepo) *ReportGenerator {
	return &ReportGenerator{outDir: outDir, reportsRepo: repo}
}

func (g *ReportGenerator) WriteUnitReport(
	ctx context.Context,
	unitGuid string,
	fileID string,
	rows []models.TSVRecord,
) (string, error) {
	if err := os.MkdirAll(g.outDir, 0o755); err != nil {
		return "", err
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].LineNo < rows[j].LineNo })

	path := filepath.Join(g.outDir, fmt.Sprintf("%s.txt", unitGuid))
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fmt.Fprintf(f, "unit_guid: %s\nfile_id: %s\ngenerated_at: %s\n\n",
		unitGuid, fileID, time.Now().Format(time.RFC3339),
	)

	for _, r := range rows {
		msg := ""
		if r.MSGID.Valid {
			msg = r.MSGID.String
		}

		fmt.Fprintf(f, "line_no=%d msg_id=%s\n", r.LineNo, msg)

		var obj any
		_ = json.Unmarshal(r.Payload, &obj)

		b, _ := json.MarshalIndent(obj, "", "  ")
		fmt.Fprintf(f, "%s\n\n", string(b))
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	rep := &models.Report{
		ID:       id.String(),
		UnitGUID: unitGuid,
		FilePath: path,
	}
	if err := g.reportsRepo.Insert(ctx, rep); err != nil {
		return "", err
	}

	return path, nil
}
