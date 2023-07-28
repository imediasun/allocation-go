package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/application"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/config"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra"
	"gitlab.hotel.tools/backend-team/allocation-go/pkg/api/openapi"
	"go.uber.org/fx"
	"time"
)

func main() {
	fxApp := fx.New(
		fx.Provide(func() context.Context { return context.Background() }),
		fx.StartTimeout(time.Second*3),
		fx.StopTimeout(time.Second*10),
		adapter.Constructors,
		infra.Constructors,

		fx.Provide(application.NewAllocatorServer),
		fx.Invoke(
			db.DB.Init,
			newAllocationServer,
		),
	)

	fxApp.Run()
}

func newAllocationServer(
	lf fx.Lifecycle,
	cfg *config.Config,
	server openapi.ServerInterface,
) {
	e := echo.New()

	openapi.RegisterHandlers(e, server)

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ = e.Close()
			return nil
		},
	})
}
