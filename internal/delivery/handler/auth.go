package handler

import (
	"avito_test/internal/domain"
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
)

type AuthService interface {
	Auth(ctx context.Context, req domain.AuthRequest) (*domain.AuthResponse, error)
}

type Auth struct {
	service AuthService
}

func NewAuth(service AuthService) Auth {
	return Auth{
		service: service,
	}
}

// Auth
// @Tags auth
// @Summary Авторизация
// @Description Авторизация пользователя
// @Accept json
// @Produce json
// @Param body domain.AuthRequest true "Данные для авторизации"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth [POST]
func (a Auth) Auth() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var req domain.AuthRequest
		if err := ctx.Bind().Body(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Errors: "invalid request body"})
		}

		res, err := a.service.Auth(ctx.Context(), req)
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return ctx.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Errors: "invalid credentials"})
		case errors.Is(err, domain.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "unauthorized"})
		case err != nil:
			return ctx.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Errors: "internal server error"})
		default:
			return ctx.Status(fiber.StatusOK).JSON(res)
		}
	}
}
