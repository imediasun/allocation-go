package http

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
)

type HTTPServer struct {
	*http.Server

	ctx    context.Context
	logger log.Logger
	cfg    *Config
}

func NewHTTPServer(ctx context.Context, logger log.Logger, cfg *Config) (*HTTPServer, error) {
	var httpServer = &http.Server{
		Addr:        net.JoinHostPort("0.0.0.0", strconv.Itoa(cfg.Port)),
		ReadTimeout: time.Second * 10,
	}

	return &HTTPServer{
		ctx:    ctx,
		logger: logger,
		cfg:    cfg,
		Server: httpServer,
	}, nil
}

func (h *HTTPServer) SetHandler(handler http.Handler) {
	h.Handler = handler
}

func (h *HTTPServer) Start() error {
	err := h.Server.ListenAndServe()
	if err != nil {
		h.logger.Error("unable to serve http server", zap.Error(err))
		return err
	}

	h.logger.Info("http server started", zap.String("addr", h.Server.Addr))

	return err
}

func (h *HTTPServer) Stop() error {
	stopCtx, cancel := context.WithTimeout(h.ctx, time.Second*5)
	defer cancel()

	err := h.Server.Shutdown(stopCtx)
	if err != nil {
		return err
	}

	h.logger.Info("http server stopped")
	return nil
}
