package routes

import (
	"github.com/Makefolder/cynero/internal/handler"
	"github.com/gofiber/fiber/v3"
)

func Setup(r fiber.Router, h *handler.Handler) {
	r.Get("/health", h.Health)

	orders := r.Group("/orders")
	orders.Get("/find-all", h.FindAll)
	orders.Post("/create", h.CreateOrder)

	webhook := r.Group("/helius-webhook")
	webhook.Post("/handle", h.HeliusWebhook)

	auth := r.Group("/auth")
	auth.Post("/login", h.AuthMerchant)
	auth.Post("/register", h.RegisterMerchant)
}
