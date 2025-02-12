package httpServer

import (
	"avito_test/internal/config"
	"avito_test/pkg/logger"
	"fmt"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

// Server struct
type Server struct {
	fiber     *fiber.App
	cfg       *config.Config
	apiLogger *logger.ApiLogger
}

func NewServer(cfg *config.Config, apiLogger *logger.ApiLogger) *Server {
	return &Server{
		fiber: fiber.New(fiber.Config{
			JSONEncoder: gojson.Marshal,
			JSONDecoder: gojson.Unmarshal,
		}),
		cfg:       cfg,
		apiLogger: apiLogger,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.fiber, s.apiLogger); err != nil {
		s.apiLogger.Fatalf("cannot map handlers: %v", err)
	}
	s.apiLogger.Infof("start server on port: %s", s.cfg.Server.Port)
	if err := s.fiber.Listen(fmt.Sprintf(":%s", s.cfg.Server.Port)); err != nil {
		s.apiLogger.Fatalf("error starting server: %v", err)
	}

	return nil
}
