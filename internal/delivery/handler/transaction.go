package handler

import (
	"avito_test/internal/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type TransactionService interface {
	Buy(ctx context.Context, userIDStr string, itemType string) error
	Send(ctx context.Context, userIDStr string, req domain.SendCoinRequest) error
	Info(ctx context.Context, userIDStr string) (*domain.InfoResponse, error)
}

type Transaction struct {
	service TransactionService
}

func NewTransaction(service TransactionService) Transaction {
	return Transaction{
		service: service,
	}
}

func (t Transaction) Buy() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userIDStr, ok := ctx.Locals("id").(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "invalid user ID format"})
		}

		itemType := ctx.Params("item")

		err := t.service.Buy(ctx.Context(), userIDStr, itemType)
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return ctx.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Errors: "invalid credentials"})
		case errors.Is(err, domain.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "unauthorized"})
		case err != nil:
			return ctx.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Errors: "internal server error"})
		default:
			return ctx.SendStatus(fiber.StatusOK)
		}
	}
}

func (t Transaction) Send() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userIDStr, ok := ctx.Locals("id").(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "invalid user ID format"})
		}

		var req domain.SendCoinRequest
		if err := ctx.Bind().Body(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Errors: "invalid request body"})
		}

		err := t.service.Send(ctx.Context(), userIDStr, req)
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return ctx.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Errors: "invalid credentials"})
		case errors.Is(err, domain.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "unauthorized"})
		case err != nil:
			return ctx.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Errors: "internal server error"})
		default:
			return ctx.SendStatus(fiber.StatusOK)
		}
	}
}

func (t Transaction) Info() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userIDStr, ok := ctx.Locals("id").(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "invalid user ID format"})
		}

		info, err := t.service.Info(ctx.Context(), userIDStr)
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return ctx.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Errors: "invalid credentials"})
		case errors.Is(err, domain.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Errors: "unauthorized"})
		case err != nil:
			return ctx.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Errors: "internal server error"})
		default:
			return ctx.Status(fiber.StatusOK).JSON(info)
		}
	}
}
