package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
)

func NewPool(ctx context.Context, logger log.Logger, config *Config, hookFactory HookFactory) (DB, error) {
	logger = logger.WithMethod(ctx, "NewPool")

	xdb, err := sqlx.Connect("mysql", config.DSN)
	if err != nil {
		logger.Error("failed to connect to db", zap.Error(err))
		return nil, err
	}

	xdb.SetMaxOpenConns(config.MaxOpenConns)
	xdb.SetMaxIdleConns(config.MaxIdleConns)
	xdb.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifetime))
	xdb.SetConnMaxIdleTime(time.Second * time.Duration(config.ConnMaxIdleTime))

	if err = xdb.Ping(); err != nil {
		logger.Error("failed to ping db", zap.Error(err))
		return nil, err
	}

	res := &mysqlDB{
		ctx:         ctx,
		config:      config,
		db:          xdb,
		logger:      logger.WithComponent(ctx, "pool"),
		hookFactory: hookFactory,
	}

	if !config.Debug {
		res.hookFactory = nil
	}

	return res, nil
}

type mysqlDB struct {
	ctx         context.Context
	config      *Config
	db          *sqlx.DB
	logger      log.Logger
	hookFactory HookFactory
}

func (d *mysqlDB) Begin(ctx context.Context) (DB, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &mysqlTx{
		ctx:         ctx,
		config:      d.config,
		tx:          tx,
		logger:      d.logger.Named("tx"),
		hookFactory: d.hookFactory,
	}, nil
}

func (d *mysqlDB) Commit() error {
	return nil
}

func (d *mysqlDB) Rollback(_ context.Context) {
}

func (d *mysqlDB) Close() error {
	d.logger.WithMethod(d.ctx, "Close").Info("db pool closed")
	return d.db.Close()
}

func (d *mysqlDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if d.hookFactory != nil {
		hook := d.hookFactory.CreateLoggingHook()
		hook.Before(ctx)
		defer hook.After(ctx, query, args...)
	}

	return d.db.ExecContext(ctx, query, args...)
}

func (d *mysqlDB) Exec(query string, args ...any) (sql.Result, error) {
	return d.ExecContext(d.ctx, query, args...)
}

func (d *mysqlDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if d.hookFactory != nil {
		hook := d.hookFactory.CreateLoggingHook()
		hook.Before(ctx)
		defer hook.After(ctx, query, args...)
	}

	return d.db.QueryContext(ctx, query, args...)
}

func (d *mysqlDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.QueryContext(d.ctx, query, args...)
}

func (d *mysqlDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if d.hookFactory != nil {
		hook := d.hookFactory.CreateLoggingHook()
		hook.Before(ctx)
		defer hook.After(ctx, query, args...)
	}

	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *mysqlDB) QueryRow(query string, args ...any) *sql.Row {
	return d.QueryRowContext(d.ctx, query, args...)
}

func (d *mysqlDB) Init() error {
	if err := d.db.Ping(); err != nil {
		return err
	}

	d.logger.Info("db pool initialized")

	return nil
}
