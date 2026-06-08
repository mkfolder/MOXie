package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/common"
	"github.com/mkfolder/moxie/internal/service"
)

type UpdateProfileRequest struct {
	Username     *string `json:"username"`
	Address      *string `json:"address"`
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

func (h *Handler) UpdatePicture(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to open file")
	}
	defer file.Close()

	url, err := h.s.UpdatePicture(c.Context(), merchantID, file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to upload file")
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}

func (h *Handler) DeletePicture(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)
	if err := h.s.DeletePicture(c.Context(), merchantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete picture")
	}

	return c.JSON(fiber.Map{"message": "picture deleted successfully"})
}
