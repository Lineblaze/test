package middleware

import (
	"avito_test/internal/jwt"
	"avito_test/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"strings"
)

type MDWManager struct {
	jwt    *jwt.Service
	logger *logger.ApiLogger
}

func NewMDWManager(jwt *jwt.Service, logger *logger.ApiLogger) *MDWManager {
	return &MDWManager{
		jwt:    jwt,
		logger: logger,
	}
}

func (mw *MDWManager) JWTMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid Authorization header",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		claims, err := mw.jwt.ParseToken(tokenParts[1])
		if err != nil {
			mw.logger.Errorf("error parsing token: %v", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		ctx.Locals("id", claims.ID)

		return ctx.Next()
	}
}
