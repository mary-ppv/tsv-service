package app

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"tsv-service/internal/app/config"
	"tsv-service/internal/app/db"
	"tsv-service/internal/app/loader"
	"tsv-service/internal/app/worker"

	"tsv-service/internal/transport/http/controllers"
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

	logger := slog.Default()

	rp := loader.InitRepositoryProvider(a.db)
	sp := loader.InitServiceProvider(ctx, rp)

	base := controllers.NewBaseController(logger, sp, rp)

	filesRepo, _ := rp.FilesRepository()
	recordsRepo, _ := rp.RecordsRepository()
	errorsRepo, _ := rp.ErrorsRepository()
	reportsRepo, _ := rp.ReportsRepository()

	reportGen := worker.NewReportGenerator(a.cfg.OutputDir, reportsRepo)
	processor := worker.NewProcessor(a.db, a.cfg.InputDir, filesRepo, recordsRepo, errorsRepo, reportGen)

	q := worker.NewQueue(256)
	q.StartWorkers(ctx, a.cfg.Workers, processor.Process)

	w := worker.NewWatcher(a.cfg.InputDir, a.cfg.PollInterval, filesRepo, q, a.db)
	go w.Run(ctx)

	unitsCtrl := admin.NewUnitsController(base)
	reportsCtrl := admin.NewReportsController(base)
	importsCtrl := admin.NewImportsController(base)

	r := router.NewRouter(a.db, unitsCtrl, reportsCtrl, importsCtrl)

	logger.Info("http listening", "addr", a.cfg.HTTPAddr)
	return r.Run(a.cfg.HTTPAddr)
}
