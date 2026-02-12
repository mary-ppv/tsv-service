package driver

import (
	"context"
	"errors"
)

var (
	ErrExecutorNotFound = errors.New("executor not found")
)

type ctxExecutor struct{}

func (ctxExecutor) String() string {
	return "driver.CtxExecutor"
}

var CtxExecutor = ctxExecutor{}

func ExecutorToContext(ctx context.Context, executor ContextExecutor) context.Context {
	return context.WithValue(ctx, CtxExecutor, executor)
}

func ExecutorFromContext(ctx context.Context) (ContextExecutor, error) {
	if ex, ok := ctx.Value(CtxExecutor).(ContextExecutor); ok {
		return ex, nil
	}

	if ex, ok := ctx.Value(CtxExecutor.String()).(ContextExecutor); ok {
		return ex, nil
	}

	return nil, ErrExecutorNotFound
}
