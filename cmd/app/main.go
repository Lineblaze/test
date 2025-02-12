package main

import (
	"avito_test/cmd/migrator"
	"avito_test/internal/config"
	"avito_test/internal/httpServer"
	"avito_test/pkg/logger"
	"log"
)

func main() {
	log.Println("starting server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("can't load config: %v", err.Error())
	}

	migrator.Migrate()

	appLogger := logger.NewApiLogger(cfg)
	err = appLogger.InitLogger()
	if err != nil {
		log.Fatalf("can't init logger: %v", err.Error())
	}

	s := httpServer.NewServer(cfg, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Errorf("server run error: %v", err)
	}
}
