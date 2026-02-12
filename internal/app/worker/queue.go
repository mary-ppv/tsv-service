package worker

import (
	"context"
	"log/slog"
	"sync"

	"tsv-service/internal/models"
)

type Queue struct {
	ch chan *models.TSVFile
}

func NewQueue(buf int) *Queue {
	return &Queue{ch: make(chan *models.TSVFile, buf)}
}

func (q *Queue) Enqueue(f *models.TSVFile) {
	q.ch <- f
}

func (q *Queue) StartWorkers(ctx context.Context, n int, fn func(context.Context, *models.TSVFile) error) {
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case f := <-q.ch:
					if f == nil {
						continue
					}
					if err := fn(ctx, f); err != nil {
						slog.Error("worker failed", "worker", workerID, "file", f.FileName, "err", err)
					}
				}
			}
		}(i + 1)
	}
}
