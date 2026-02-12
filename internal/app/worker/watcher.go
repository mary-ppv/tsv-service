package worker

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tsv-service/internal/repository"
)

type Watcher struct {
	inputDir     string
	pollInterval time.Duration
	filesRepo    *repository.FilesRepo
	queue        *Queue
}

func NewWatcher(inputDir string, poll time.Duration, filesRepo *repository.FilesRepo, q *Queue) *Watcher {
	return &Watcher{inputDir: inputDir, pollInterval: poll, filesRepo: filesRepo, queue: q}
}

func (w *Watcher) Run(ctx context.Context) {
	slog.Info("watcher started")
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
	slog.Info("scan input dir", "dir", w.inputDir)
	entries, err := os.ReadDir(w.inputDir)
	if err != nil {
		slog.Error("read input dir failed", "err", err)
		return
	}
	slog.Info("input entries", "count", len(entries))

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		slog.Info("found file", "name", name)

		if !strings.HasSuffix(strings.ToLower(name), ".tsv") {
			slog.Info("skip non-tsv", "name", name)
			continue
		}

		slog.Info("tsv accepted", "name", name)

		full := filepath.Join(w.inputDir, name)
		row, created, err := w.filesRepo.EnsureQueued(ctx, full)
		if err != nil {
			slog.Error("ensure queued failed", "file", name, "err", err)
			continue
		}

		slog.Info("ensure queued result",
			"file", name,
			"created", created,
			"row_id", row.ID,
			"row_file", row.FileName,
			"row_sha", row.FileSha256,
			"row_status", row.Status,
		)
		
		if created {
			w.queue.Enqueue(row)
			slog.Info("queued file", "file", name)
		}
	}
}
