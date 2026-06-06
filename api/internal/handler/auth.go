package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/auth"
	"github.com/mkfolder/moxie/internal/common"
	"github.com/mkfolder/moxie/internal/service"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TwoFactorRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"totp_code"`
}

func (h *Handler) AuthMerchant(c fiber.Ctx) error {
	var body AuthRequest

	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidateEmail(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.s.AuthMerchant(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if result.NeedTwoFA {
		return c.JSON(fiber.Map{"need_2fa": true})
	}

	h.setAuthCookies(c, result.AuthResponse)

	return c.JSON(fiber.Map{"merchant": result.AuthResponse.Merchant})
}

func (h *Handler) RegisterMerchant(c fiber.Ctx) error {
	var body AuthRequest
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidateEmail(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := common.ValidatePassword(body.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	merchant, err := h.s.RegisterMerchant(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"merchant_id": merchant.ID, "email": merchant.Email})
}

func (h *Handler) RefreshMerchant(c fiber.Ctx) error {
	refreshTokenStr := c.Cookies("refresh_token")
	if refreshTokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
	}

	claims, err := auth.ValidateRefreshToken(h.authCfg.JWTSecret, refreshTokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired refresh token"})
	}

	merchant, err := h.s.FindMerchantByID(c.Context(), claims.MerchantID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "merchant not found"})
	}

	res, err := h.s.GenerateTokens(c.Context(), merchant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.setAuthCookies(c, res)

	return c.JSON(fiber.Map{"merchant": merchant})
}

func (h *Handler) Me(c fiber.Ctx) error {
	merchantID := c.Locals("merchant_id").(uuid.UUID)
	merchant, err := h.s.FindMerchantByID(c.Context(), merchantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "merchant not found"})
	}

	isServiceEnabled := merchant.HeliusAPIKey != nil && merchant.WebhookURL != nil && merchant.Address != nil
	merchant.IsServiceEnabled = isServiceEnabled
	return c.JSON(merchant)
}

func (h *Handler) LogoutMerchant(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "strict",
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) SetupTwoFactor(c fiber.Ctx) error {
	var body AuthRequest

	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	setup, err := h.s.SetupTOTP(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(setup)
}

func (h *Handler) VerifyTwoFactor(c fiber.Ctx) error {
	var body TwoFactorRequest

	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.s.VerifyTOTPAndAuth(c.Context(), body.Email, body.Password, body.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	h.setAuthCookies(c, res)

	return c.JSON(fiber.Map{"merchant": res.Merchant})
}

func (h *Handler) setAuthCookies(c fiber.Ctx, res *service.AuthResponse) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Path:     "/",
		MaxAge:   int(h.authCfg.AccessTokenTTL.Seconds()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Path:     "/",
		MaxAge:   int(h.authCfg.RefreshTokenTTL.Seconds()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "strict",
	})
}
