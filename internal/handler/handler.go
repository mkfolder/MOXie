package handler

import (
	"time"

	"github.com/Makefolder/moxie/internal/service"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	s         *service.Service
	startedAt time.Time
}

func New(s *service.Service) *Handler {
	return &Handler{
		s:         s,
		startedAt: time.Now().UTC(),
	}
}

func (h *Handler) Health(c fiber.Ctx) error {
	dbStatus := "healthy"

	now := time.Now().UTC()
	uptime := time.Since(h.startedAt).String()

	if err := h.s.PingDB(c.Context()); err != nil {
		dbStatus = "unhealthy"
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"now":    now.Format(time.DateTime),
		"uptime": uptime,
		"db":     dbStatus,
	})
}
