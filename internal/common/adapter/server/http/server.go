package http

import (
	"context"
	"net/http"
	"sort"

	"go.uber.org/zap"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
)

type Server interface {
	Start() error
	Stop()
}

type server struct {
	ctx        context.Context
	logger     log.Logger
	httpServer *HTTPServer
	serveMux   *http.ServeMux
}

func NewHTTP(
	ctx context.Context,
	logger log.Logger,

	httpServer *http.Server,
	serveMux *http.ServeMux,
	middlewares []Middleware,
) Server {
	sort.Sort(ByOrder(middlewares))

	httpServer.Handler = serveMux

	return &server{
		ctx:    ctx,
		logger: logger,
		//	httpServer: httpServer,
		serveMux: serveMux,
	}
}

func (s *server) Start() error {
	errs := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.logger.Error("serve http server", zap.Error(err))
			errs <- err
		}
	}()

	return <-errs
}

func (s *server) Stop() {
	if err := s.httpServer.Close(); err != nil {
		s.logger.Error("stop http server", zap.Error(err))
	}
}
