package adapter

import (
	"context"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/config"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	config.NewConfig,
	db.NewHookFactory,
	NewFxLogger,
	NewFxMysqlPool,
)

func NewFxLogger(cfg *config.Config) (log.Logger, error) {
	return log.NewLogger(cfg.Logger)
}

func NewFxMysqlPool(lf fx.Lifecycle, ctx context.Context, logger log.Logger, cfg *config.Config, hookFactory db.HookFactory) (db.DB, error) {
	pool, err := db.NewPool(ctx, logger, cfg.DB, hookFactory)
	lf.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pool.Close()
		},
	})

	return pool, err
}
