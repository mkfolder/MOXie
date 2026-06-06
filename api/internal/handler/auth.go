package handler

import (
	"encoding/json"
	"net/url"

	"github.com/gofiber/fiber/v3"
	"github.com/mkfolder/moxie/internal/common"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	AuthRequest
	Address    string `json:"address"`
	WebhookURL string `json:"webhook_url"`
}

func (h *Handler) AuthMerchant(c fiber.Ctx) error {
	var body AuthRequest

	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidateEmail(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.s.AuthMerchant(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *Handler) RegisterMerchant(c fiber.Ctx) error {
	var body RegisterRequest

	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidateEmail(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidatePassword(body.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Address == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}

	webhookURL, err := url.Parse(body.WebhookURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	merchant, err := h.s.RegisterMerchant(
		c.Context(), body.Email, body.Password, body.Address, webhookURL,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(merchant)
}
