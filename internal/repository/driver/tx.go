package driver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/friendsofgo/errors"
)

func TXExecutor(ctx context.Context, executor ContextExecutor) (ContextExecutor, error) {
	switch executor.(type) {
	case *sql.Tx:
		return executor, nil
	default:
		beginer, ok := executor.(boil.ContextBeginner)
		if !ok {
			return nil, fmt.Errorf("unsupported executor type %T", executor)
		}

		return beginer.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	}
}

func BeginTxInContext(ctx context.Context) (context.Context, error) {
	executor, err := ExecutorFromContext(ctx)
	if err != nil {
		return nil, err
	}

	executor, err = TXExecutor(ctx, executor)
	if err != nil {
		return nil, err
	}

	return ExecutorToContext(ctx, executor), nil
}

func EndTxInContext(ctx context.Context, needRollback bool) error {
	executor, err := ExecutorFromContext(ctx)
	if err != nil {
		return err
	}

	if txExecutor, ok := executor.(*sql.Tx); ok {
		if needRollback {
			return txExecutor.Rollback()
		}

		return txExecutor.Commit()
	}

	return nil
}

func ExecuteTransaction(ctx context.Context, dbAction func(context.Context, ContextExecutor) error) (err error) {
	ctx, err = BeginTxInContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin tx")
	}

	defer func() {
		needRollback := err != nil
		if txErr := errors.Wrap(EndTxInContext(ctx, needRollback), "failed to end tx"); txErr != nil {
			if err != nil {
				err = errors.Wrap(txErr, fmt.Sprintf("with err: %v", err))
				return
			}

			err = txErr
		}
	}()

	executor, err := ExecutorFromContext(ctx)
	if err != nil {
		return err
	}

	err = dbAction(ctx, executor)

	return err
}
