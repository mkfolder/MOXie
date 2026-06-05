package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Makefolder/moxie/internal/models"
	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	Token    string           `json:"token"`
	Merchant *models.Merchant `json:"merchant"`
}

func (s *Service) RegisterMerchant(
	ctx context.Context,
	email, password, address string,
	webhookURL *url.URL,
) (*models.Merchant, error) {
	var merchant models.Merchant

	merchant.Address = address
	merchant.Email = email
	merchant.WebhookURL = new(string)
	*merchant.WebhookURL = webhookURL.String()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	merchant.PasswdHash = hash

	if err := s.merchants.Create(ctx, &merchant); err != nil {
		return nil, fmt.Errorf("failed to register new merchant: %w", err)
	}

	go s.hc.CreateWebhook(ctx, []string{address})
	return &merchant, nil
}

func (s *Service) AuthMerchant(ctx context.Context, email, password string) (*AuthResponse, error) {
	var merchant models.Merchant

	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&merchant).Error; err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(merchant.PasswdHash, []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// todo!: JWT token issuing
	token := base58.Encode(merchant.ID[:])
	return &AuthResponse{Token: token, Merchant: &merchant}, nil
}
