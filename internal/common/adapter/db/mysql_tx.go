package db

import (
	"context"

	"database/sql"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
)

type mysqlTx struct {
	ctx         context.Context
	config      *Config
	tx          *sql.Tx
	logger      log.Logger
	hookFactory HookFactory
}

func (t *mysqlTx) Begin(context.Context) (DB, error) {
	return t, nil
}

func (t *mysqlTx) Commit() error {
	return t.tx.Commit()
}

func (t *mysqlTx) Rollback(ctx context.Context) {
	if err := t.tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		t.logger.WithMethod(ctx, "Rollback").Error("failed to rollback tx", zap.Error(err))
	}
}

func (t *mysqlTx) Close() error {
	return nil
}

func (t *mysqlTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if t.hookFactory != nil {
		hook := t.hookFactory.CreateLoggingHook()
		hook.Before(ctx)
		defer hook.After(ctx, query, args...)
	}

	return t.tx.ExecContext(ctx, query, args...)
}

func (t *mysqlTx) Exec(query string, args ...any) (sql.Result, error) {
	return t.ExecContext(t.ctx, query, args...)
}

func (t *mysqlTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if t.hookFactory != nil {
		hook := t.hookFactory.CreateLoggingHook()
		hook.Before(ctx)
		defer hook.After(ctx, query, args...)
	}

	return t.tx.QueryContext(ctx, query, args...)
}

func (t *mysqlTx) Query(query string, args ...any) (*sql.Rows, error) {
	return t.QueryContext(t.ctx, query, args...)
}

func (t *mysqlTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if t.hookFactory != nil {
		hook := t.hookFactory.CreateLoggingHook()
		hook.Before(ctx)
		defer hook.After(ctx, query, args...)
	}

	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *mysqlTx) QueryRow(query string, args ...any) *sql.Row {
	return t.QueryRowContext(t.ctx, query, args...)
}

func (t *mysqlTx) Init() error {
	return nil
}
