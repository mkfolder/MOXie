package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/constants"
)

func (h *Handler) FindAll(c fiber.Ctx) error {
	orders, err := h.s.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (h *Handler) FindOrder(c fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order, err := h.s.FindOrder(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}

func (h *Handler) CreateOrder(c fiber.Ctx) error {
	var requestBody struct {
		MerchantID string          `json:"merchant_id"`
		RawAmount  uint64          `json:"raw_amount"`
		CustomData json.RawMessage `json:"custom_data"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	merchantID, err := uuid.Parse(requestBody.MerchantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order, err := h.s.CreateOrder(
		c.Context(),
		merchantID,
		requestBody.RawAmount,
		requestBody.CustomData,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"order":       order,
		"qrcode_data": fmt.Sprintf("solana:https://%s/solpay/%s", constants.Subdomain, order.ID.String()),
	})
}
