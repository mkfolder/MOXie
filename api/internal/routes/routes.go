package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mkfolder/moxie/internal/handler"
	"github.com/mkfolder/moxie/internal/middleware"
)

func Setup(r fiber.Router, h *handler.Handler) {
	api := r.Group("/api")

	api.Get("/health", h.Health)
	api.Post("/service/orders/create", h.CreateOrder)

	auth := api.Group("/auth")
	auth.Post("/login", h.AuthMerchant)
	auth.Post("/register", h.RegisterMerchant)
	auth.Post("/refresh", h.RefreshMerchant)
	auth.Post("/logout", h.LogoutMerchant)
	auth.Post("/2fa/setup", h.SetupTwoFactor)
	auth.Post("/2fa/verify", h.VerifyTwoFactor)

	protected := api.Group("", middleware.AuthRequired(h.AuthCfg().JWTSecret))

	protected.Get("/auth/me", h.Me)

	profile := protected.Group("/profile")
	profile.Put("/update", h.UpdateProfile)
	profile.Put("/update-picture", h.UpdatePicture)
	profile.Delete("/delete-picture", h.DeletePicture)
	profile.Put("/password", h.ChangePassword)
	profile.Get("/helius-key", h.GetHeliusKey)

	orders := protected.Group("/orders")
	orders.Get("/find-all", h.FindAll)
	orders.Get("/find/:id", h.FindOrder)

	solpay := protected.Group("/solpay")
	solpay.Get("/:id", h.GetSolPayMetadata)
	solpay.Post("/:id", h.BuildSolPayTransaction)

	webhook := api.Group("/helius-webhook")
	webhook.Post("/handle", h.HeliusWebhook)
}
