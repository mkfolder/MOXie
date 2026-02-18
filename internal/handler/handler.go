package handler

import (
	"github.com/Makefolder/cynero/internal/service"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	s *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) Test(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": h.s.Test(),
	})
}
