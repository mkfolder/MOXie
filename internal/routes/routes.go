package routes

import (
	"github.com/Makefolder/cynero/internal/handler"
	"github.com/gofiber/fiber/v3"
)

func Setup(r fiber.Router, h *handler.Handler) {
	r.Get("/test", h.Test)
}
