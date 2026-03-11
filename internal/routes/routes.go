package routes

import (
	"github.com/Makefolder/cynero/internal/handler"
	"github.com/gofiber/fiber/v3"
)

func Setup(r fiber.Router, h *handler.Handler) {
	r.Get("/health", h.Health)
	r.Get("/find-all", h.FindAll)
	r.Post("/create-order", h.CreateOrder)
	r.Post("/handle/helius-webhook", h.HeliusWebhook)
	r.Post("/subscribe/helius-webhook", h.SubscribeHeliusWebhook)
}
