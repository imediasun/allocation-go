package log_test

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
)

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerSuite))
}

type LoggerSuite struct {
	suite.Suite
}

func (s *LoggerSuite) TestAll() {
	logger := log.MustDefaultJSONLogger("")
	s.NotEmpty(logger)
	logger = logger.WithCtx(context.TODO())
	s.NotEmpty(logger)
	logger = logger.Named(gofakeit.Name())
	s.NotEmpty(logger)

	logger, err := log.NewLogger(&zap.Config{})
	s.Empty(logger)
	s.Contains(err.Error(), "no encoder name specified")
}

func (s *LoggerSuite) TestConsoleLogger() {
	s.NotPanics(func() { log.MustDefaultConsoleLogger("trace") })
}

func (s *LoggerSuite) TestJsonLogger() {
	s.NotPanics(func() { log.MustDefaultJSONLogger("trace") })
}

func (s *LoggerSuite) TestNamed() {
	logger := log.MustDefaultConsoleLogger("trace")
	s.NotEmpty(logger)

	logger = logger.Named("foo")
	logger.Debug("foo message")
}
