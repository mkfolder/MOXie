package service

import (
	"context"
	"fmt"

	"github.com/mkfolder/moxie/internal/auth"
	"github.com/mkfolder/moxie/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	Merchant     *models.Merchant `json:"merchant"`
}

type LoginResult struct {
	AuthResponse *AuthResponse
	NeedTwoFA    bool
}

func (s *Service) RegisterMerchant(
	ctx context.Context,
	email, password string,
) (*models.Merchant, error) {
	var merchant models.Merchant

	merchant.Email = email
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	merchant.PasswdHash = hash
	if err := s.merchants.Create(ctx, &merchant); err != nil {
		return nil, fmt.Errorf("failed to register new merchant: %w", err)
	}

	isServiceEnabled := merchant.HeliusAPIKey != nil && merchant.WebhookURL != nil && merchant.Address != nil
	merchant.IsServiceEnabled = isServiceEnabled

	return &merchant, nil
}

func (s *Service) AuthMerchant(ctx context.Context, email, password string) (*LoginResult, error) {
	var merchant models.Merchant

	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&merchant).Error; err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(merchant.PasswdHash, []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	if merchant.TOTPEnabled {
		return &LoginResult{NeedTwoFA: true}, nil
	}

	isServiceEnabled := merchant.HeliusAPIKey != nil && merchant.WebhookURL != nil && merchant.Address != nil
	merchant.IsServiceEnabled = isServiceEnabled

	resp, err := s.GenerateTokens(ctx, &merchant)
	if err != nil {
		return nil, err
	}

	return &LoginResult{AuthResponse: resp}, nil
}

func (s *Service) GenerateTokens(ctx context.Context, merchant *models.Merchant) (*AuthResponse, error) {
	accessToken, err := auth.GenerateAccessToken(
		s.authCfg.JWTSecret, merchant.ID, merchant.Email, s.authCfg.AccessTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(
		s.authCfg.JWTSecret, merchant.ID, s.authCfg.RefreshTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Merchant:     merchant,
	}, nil
}

func (s *Service) SetupTOTP(ctx context.Context, email, password string) (*auth.TOTPSetup, error) {
	var merchant models.Merchant
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&merchant).Error; err != nil {
		return nil, fmt.Errorf("merchant not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(merchant.PasswdHash, []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	setup, err := auth.GenerateTOTPSecret(email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate totp secret: %w", err)
	}

	merchant.TOTPSecret = setup.Secret
	if err := s.db.WithContext(ctx).Save(&merchant).Error; err != nil {
		return nil, fmt.Errorf("failed to save totp secret: %w", err)
	}

	return setup, nil
}

func (s *Service) VerifyTOTPAndAuth(ctx context.Context, email, password, code string) (*AuthResponse, error) {
	var merchant models.Merchant
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&merchant).Error; err != nil {
		return nil, fmt.Errorf("merchant not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(merchant.PasswdHash, []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	if merchant.TOTPSecret == "" {
		return nil, fmt.Errorf("totp not set up")
	}

	if !auth.ValidateTOTPCode(merchant.TOTPSecret, code) {
		return nil, fmt.Errorf("invalid totp code")
	}

	if !merchant.TOTPEnabled {
		merchant.TOTPEnabled = true
		if err := s.db.WithContext(ctx).Save(&merchant).Error; err != nil {
			return nil, fmt.Errorf("failed to enable totp: %w", err)
		}
	}

	isServiceEnabled := merchant.HeliusAPIKey != nil && merchant.WebhookURL != nil && merchant.Address != nil
	merchant.IsServiceEnabled = isServiceEnabled

	return s.GenerateTokens(ctx, &merchant)
}
