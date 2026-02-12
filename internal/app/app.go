package app

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"tsv-service/internal/services/imports"
	"tsv-service/internal/services/reports"
	"tsv-service/internal/services/units"

	"tsv-service/internal/app/config"
	"tsv-service/internal/app/db"
	"tsv-service/internal/app/worker"
	"tsv-service/internal/repository"
	admin "tsv-service/internal/transport/http/controllers/admin"
	"tsv-service/internal/transport/http/router"
)

type App struct {
	cfg config.Config
	db  *sql.DB
}

func New(cfg config.Config) (*App, error) {
	conn, err := db.Open(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &App{cfg: cfg, db: conn}, nil
}

func (a *App) Run(ctx context.Context) error {
	_ = os.MkdirAll(a.cfg.InputDir, 0o755)
	_ = os.MkdirAll(a.cfg.OutputDir, 0o755)

	filesRepo := repository.NewFilesRepo(a.db)
	recordsRepo := repository.NewRecordsRepo(a.db)
	errorsRepo := repository.NewErrorsRepo(a.db)
	reportsRepo := repository.NewReportsRepo(a.db)
	unitsRepo := repository.NewUnitsRepo(a.db)

	reportGen := worker.NewReportGenerator(a.cfg.OutputDir, reportsRepo)
	processor := worker.NewProcessor(a.db, a.cfg.InputDir, filesRepo, recordsRepo, errorsRepo, reportGen)

	q := worker.NewQueue(256)
	q.StartWorkers(ctx, a.cfg.Workers, processor.Process)

	w := worker.NewWatcher(a.cfg.InputDir, a.cfg.PollInterval, filesRepo, q)
	go w.Run(ctx)

	unitsSvc := units.NewService(unitsRepo, recordsRepo)
	reportsSvc := reports.NewService(reportsRepo)
	importsSvc := imports.NewService(unitsRepo)

	unitsCtrl := admin.NewUnitsController(unitsSvc)
	reportsCtrl := admin.NewReportsController(reportsSvc)
	importsCtrl := admin.NewImportsController(importsSvc)

	r := router.NewRouter(unitsCtrl, reportsCtrl, importsCtrl)

	slog.Info("http listening", "addr", a.cfg.HTTPAddr)
	return r.Run(a.cfg.HTTPAddr)
}
