package handler

import (
	"encoding/json"

	"github.com/Makefolder/cynero/internal/helius"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) HeliusWebhook(c fiber.Ctx) error {
	var transacitons []helius.Transaction
	if err := json.Unmarshal(c.Body(), &transacitons); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	go h.s.HandleWebhook(c.Context(), transacitons)
	return c.JSON(fiber.Map{"status": "ok"})
}
