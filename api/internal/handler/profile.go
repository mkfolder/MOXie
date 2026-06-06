package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/service"
	"github.com/mkfolder/moxie/internal/common"
)

type UpdateProfileRequest struct {
	Username     *string `json:"username"`
	Address      *string `json:"address"`
	AvatarURL    *string `json:"avatar_url"`
	WebhookURL   *string `json:"webhook_url"`
	HeliusAPIKey *string `json:"helius_api_key"`
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)

	var body UpdateProfileRequest
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	svcReq := &service.UpdateProfileRequest{
		Username:     body.Username,
		Address:      body.Address,
		AvatarURL:    body.AvatarURL,
		WebhookURL:   body.WebhookURL,
		HeliusAPIKey: body.HeliusAPIKey,
	}

	merchant, err := h.s.UpdateMerchantProfile(c.Context(), merchantID, svcReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(merchant)
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *Handler) ChangePassword(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)

	var body ChangePasswordRequest
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidatePassword(body.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	svcReq := &service.ChangePasswordRequest{
		CurrentPassword: body.CurrentPassword,
		NewPassword:     body.NewPassword,
	}

	if err := h.s.ChangePassword(c.Context(), merchantID, svcReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "password updated successfully"})
}

func (h *Handler) GetHeliusKey(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)

	key, err := h.s.GetHeliusAPIKey(c.Context(), merchantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"helius_api_key": key})
}
