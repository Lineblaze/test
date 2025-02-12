package httpServer

import (
	"avito_test/internal/delivery/handler"
	"avito_test/internal/delivery/routes"
	"avito_test/internal/jwt"
	"avito_test/internal/middleware"
	"avito_test/internal/repository"
	"avito_test/internal/service"
	"avito_test/pkg/logger"
	storage "avito_test/pkg/storage/postgres"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	serverLogger "github.com/gofiber/fiber/v3/middleware/logger"
)

func (s *Server) MapHandlers(app *fiber.App, logger *logger.ApiLogger) error {
	db, err := storage.InitPsqlDB(s.cfg)
	if err != nil {
		logger.Fatalf("failed to initialize PostgreSQL DB: %v", err)
	}

	jwtService := jwt.NewJWTService(s.cfg)

	authRepo := repository.NewAuth(db)
	authService := service.NewAuth(authRepo, jwtService)
	authHandler := handler.NewAuth(authService)

	transactionRepo := repository.NewTransaction(db)
	transactionService := service.NewTransaction(transactionRepo)
	transactionHandler := handler.NewTransaction(transactionService)

	app.Use(serverLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{},
		AllowHeaders: []string{},
	}))

	mw := middleware.NewMDWManager(jwtService, logger)

	authGroup := app.Group("/api")
	transactionGroup := app.Group("/api/transaction/")
	transactionGroup.Use(mw.JWTMiddleware())
	routes.MapAuthRoutes(authGroup, authHandler)
	routes.MapTransactionRoutes(transactionGroup, transactionHandler)

	return nil
}
