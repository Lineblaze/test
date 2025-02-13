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

// Buy
// @Tags transactions
// @Summary Покупка предмета
// @Description Совершение покупки предмета пользователем
// @Accept json
// @Produce json
// @Param item path string true "Тип предмета"
// @Success 200 "Успешная покупка"
// @Failure 400 {object} domain.ErrorResponse "Некорректные учетные данные"
// @Failure 401 {object} domain.ErrorResponse "Неавторизованный доступ"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /transactions/buy/{item} [GET]
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

// Send
// @Tags transactions
// @Summary Отправка монет
// @Description Отправка монет другому пользователю
// @Accept json
// @Produce json
// @Param body body domain.SendCoinRequest true "Данные для перевода"
// @Success 200 "Перевод успешно выполнен"
// @Failure 400 {object} domain.ErrorResponse "Некорректные учетные данные или тело запроса"
// @Failure 401 {object} domain.ErrorResponse "Неавторизованный доступ"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /transactions/send [POST]
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

// Info
// @Tags transactions
// @Summary Информация о транзакциях
// @Description Получение информации о транзакциях пользователя
// @Accept json
// @Produce json
// @Success 200 {object} domain.InfoResponse "Информация о транзакциях"
// @Failure 400 {object} domain.ErrorResponse "Некорректные учетные данные"
// @Failure 401 {object} domain.ErrorResponse "Неавторизованный доступ"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /transactions/info [GET]
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
