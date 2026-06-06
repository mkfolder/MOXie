package handler

import (
	"encoding/base64"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/constants"
)

type SolPayMetadata struct {
	Label string `json:"label"`
	Icon  string `json:"icon"`
}

type SolPayTransaction struct {
	Transaction string `json:"transaction"`
	Message     string `json:"message"`
}

type BuildSolPayTransactionRequest struct {
	Address string `json:"address"`
}

func (h *Handler) GetSolPayMetadata(c fiber.Ctx) error {
	icon := fmt.Sprintf("https://%s/public/icon.png", constants.Domain)
	return c.JSON(SolPayMetadata{
		Label: "MOXie Payment Processor",
		Icon:  icon,
	})
}

func (h *Handler) BuildSolPayTransaction(c fiber.Ctx) error {
	var req BuildSolPayTransactionRequest
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	sender, err := solana.PublicKeyFromBase58(req.Address)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid sender address",
		})
	}

	txn, err := h.s.BuildOrderTransaction(c.Context(), sender, orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	serialized, err := txn.MarshalBinary()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to serialize transaction")
	}

	return c.JSON(SolPayTransaction{
		Transaction: base64.StdEncoding.EncodeToString(serialized),
		Message:     fmt.Sprintf("Transaction for %s", orderIDStr),
	})
}
