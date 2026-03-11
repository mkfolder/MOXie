package handler

import (
	"encoding/json"
	"time"

	"github.com/Makefolder/cynero/internal/service"
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

func (h *Handler) FindAll(c fiber.Ctx) error {
	orders, err := h.s.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (h *Handler) CreateOrder(c fiber.Ctx) error {
	var requestBody struct {
		Address    string          `json:"address"`
		RawAmount  uint64          `json:"raw_amount"`
		CustomData json.RawMessage `json:"custom_data"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	payment, err := h.s.CreateOrder(
		c.Context(),
		requestBody.Address,
		requestBody.RawAmount,
		requestBody.CustomData,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(payment)
}

func (h *Handler) HeliusWebhook(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
