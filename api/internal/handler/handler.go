package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mkfolder/moxie/internal/config"
	"github.com/mkfolder/moxie/internal/service"
)

type Handler struct {
	s         *service.Service
	startedAt time.Time
	authCfg   config.Auth
}

func New(s *service.Service, authCfg config.Auth) *Handler {
	return &Handler{
		s:         s,
		startedAt: time.Now().UTC(),
		authCfg:   authCfg,
	}
}

func (h *Handler) AuthCfg() config.Auth {
	return h.authCfg
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
