package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/constants"
	"github.com/mkfolder/moxie/internal/models"
)

type OrderWithQRCode struct {
	models.Order
	QRCode string `json:"qrcode_data"`
}

func (h *Handler) FindAll(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	orders, total, err := h.s.FindAll(c.Context(), merchantID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var ordersWithQRCode []OrderWithQRCode
	for _, order := range orders {
		ordersWithQRCode = append(ordersWithQRCode, OrderWithQRCode{
			Order:  order,
			QRCode: getQRCode(order.ID),
		})
	}

	return c.JSON(fiber.Map{
		"data":   ordersWithQRCode,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
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
		"qrcode_data": getQRCode(order.ID),
	})
}

func getQRCode(orderID uuid.UUID) string {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = constants.DomainFallback
	}
	return fmt.Sprintf("solana:https://%s/api/solpay/%s", domain, orderID.String())
}
