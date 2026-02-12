package worker

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tsv-service/internal/repository"
	"tsv-service/internal/repository/driver"
)

type Watcher struct {
	inputDir     string
	pollInterval time.Duration
	filesRepo    *repository.FilesRepo
	queue        *Queue
	exec         driver.ContextExecutor
}

func NewWatcher(
	inputDir string,
	poll time.Duration,
	filesRepo *repository.FilesRepo,
	q *Queue,
	exec driver.ContextExecutor,
) *Watcher {
	return &Watcher{
		inputDir:     inputDir,
		pollInterval: poll,
		filesRepo:    filesRepo,
		queue:        q,
		exec:         exec,
	}
}

func (w *Watcher) Run(ctx context.Context) {
	slog.Info("watcher started")

	ctx = driver.ExecutorToContext(ctx, w.exec)

	t := time.NewTicker(w.pollInterval)
	defer t.Stop()

	w.scanOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.scanOnce(ctx)
		}
	}
}

func (w *Watcher) scanOnce(ctx context.Context) {
	entries, err := os.ReadDir(w.inputDir)
	if err != nil {
		slog.Error("read input dir failed", "dir", w.inputDir, "err", err)
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".tsv") {
			continue
		}

		full := filepath.Join(w.inputDir, name)

		row, created, err := w.filesRepo.EnsureQueued(ctx, full)
		if err != nil {
			slog.Error("ensure queued failed", "file", name, "err", err)
			continue
		}

		if created {
			w.queue.Enqueue(row)
			slog.Info("queued file", "file", name)
		}
	}
}
