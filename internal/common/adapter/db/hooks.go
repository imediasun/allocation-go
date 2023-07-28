package db

import (
	"context"
	"time"

	"go.uber.org/zap"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
)

type Hook interface {
	Before(context.Context)
	After(ctx context.Context, query string, args ...any)
}

type loggingHook struct {
	logger log.Logger
	start  time.Time
}

func (h *loggingHook) Before(ctx context.Context) {
	h.start = time.Now()
}

func (h *loggingHook) After(ctx context.Context, query string, args ...any) {
	h.logger.WithMethod(ctx, "After").Debug(
		"exec",
		zap.String("query", query),
		zap.Any("args", args),
		zap.String("duration", time.Since(h.start).String()),
	)
}

type HookFactory interface {
	CreateLoggingHook() Hook
}

type hookFactory struct {
	logger log.Logger
}

func (f *hookFactory) CreateLoggingHook() Hook {
	return &loggingHook{
		logger: f.logger.Named("logging_hook"),
	}
}

func NewHookFactory(logger log.Logger) HookFactory {
	return &hookFactory{
		logger: logger,
	}
}
