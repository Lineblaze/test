package routes

import "github.com/gofiber/fiber/v3"

type AuthHandler interface {
	Auth() fiber.Handler
}

type TransactionHandler interface {
	Buy() fiber.Handler
	Send() fiber.Handler
	Info() fiber.Handler
}

func MapAuthRoutes(r fiber.Router, h AuthHandler) {
	r.Post(`/auth`, h.Auth())
}

func MapTransactionRoutes(r fiber.Router, h TransactionHandler) {
	r.Get(`/info`, h.Info())
	r.Get(`/buy/:item`, h.Buy())
	r.Post(`/sendCoin`, h.Send())
}
