package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/common"
	"github.com/mkfolder/moxie/internal/crypto"
	"github.com/mkfolder/moxie/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UpdateProfileRequest struct {
	Username     *string `json:"username"`
	Address      *string `json:"address"`
	AvatarURL    *string `json:"avatar_url"`
	WebhookURL   *string `json:"webhook_url"`
	HeliusAPIKey *string `json:"helius_api_key"`
}

func (s *Service) UpdateMerchantProfile(ctx context.Context, merchantID uuid.UUID, req *UpdateProfileRequest) (*models.Merchant, error) {
	merchant, err := s.FindMerchantByID(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		if *req.Username == "" {
			return nil, errors.New("username cannot be empty")
		}
		merchant.Username = *req.Username
	}

	if req.Address != nil {
		if *req.Address == "" {
			merchant.Address = nil
		} else {
			merchant.Address = req.Address
		}
	}

	if req.AvatarURL != nil {
		if *req.AvatarURL == "" {
			merchant.AvatarURL = nil
		} else {
			merchant.AvatarURL = req.AvatarURL
		}
	}

	if req.WebhookURL != nil {
		if *req.WebhookURL == "" {
			merchant.WebhookURL = nil
		} else {
			merchant.WebhookURL = req.WebhookURL
		}
	}

	if req.HeliusAPIKey != nil {
		if *req.HeliusAPIKey == "" {
			merchant.HeliusAPIKey = nil
		} else {
			encrypted, err := crypto.Encrypt(*req.HeliusAPIKey, s.authCfg.EncryptionKey)
			if err != nil {
				return nil, err
			}
			merchant.HeliusAPIKey = &encrypted
		}
	}

	if err := s.merchants.Update(ctx, merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (s *Service) ChangePassword(ctx context.Context, merchantID uuid.UUID, req *ChangePasswordRequest) error {
	if err := common.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	merchant, err := s.FindMerchantByID(ctx, merchantID)
	if err != nil {
		return fmt.Errorf("merchant not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(merchant.PasswdHash, []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	merchant.PasswdHash = hash
	if err := s.merchants.Update(ctx, merchant); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *Service) GetHeliusAPIKey(ctx context.Context, merchantID uuid.UUID) (*string, error) {
	merchant, err := s.FindMerchantByID(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	if merchant.HeliusAPIKey == nil {
		return nil, nil
	}

	decrypted, err := crypto.Decrypt(*merchant.HeliusAPIKey, s.authCfg.EncryptionKey)
	if err != nil {
		return nil, err
	}

	return &decrypted, nil
}
