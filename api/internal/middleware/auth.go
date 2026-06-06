package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mkfolder/moxie/internal/auth"
)

func AuthRequired(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenStr := c.Cookies("access_token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing access token"})
		}

		claims, err := auth.ValidateAccessToken(jwtSecret, tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired access token"})
		}

		c.Locals("merchant_id", claims.MerchantID)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}
