package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/common"
	"github.com/mkfolder/moxie/internal/crypto"
	"github.com/mkfolder/moxie/internal/models"
	"github.com/mkfolder/moxie/pkg/s3client"
	"golang.org/x/crypto/bcrypt"
)

type UpdateProfileRequest struct {
	Username     *string `json:"username"`
	Address      *string `json:"address"`
	WebhookURL   *string `json:"webhook_url"`
	HeliusAPIKey *string `json:"helius_api_key"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

var allowedMIMEs = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/webp": true,
}

var MIMEExtensions = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpeg",
	"image/webp": "webp",
}

func (s *Service) UpdateMerchantProfile(
	ctx context.Context,
	merchantID uuid.UUID,
	req *UpdateProfileRequest,
) (*models.Merchant, error) {
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

	isServiceEnabled := merchant.HeliusAPIKey != nil && merchant.WebhookURL != nil && merchant.Address != nil
	merchant.IsServiceEnabled = isServiceEnabled

	return merchant, nil
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

func (s *Service) UpdatePicture(ctx context.Context, merchantID uuid.UUID, file multipart.File) (string, error) {
	merchant, err := s.FindMerchantByID(ctx, merchantID)
	if err != nil {
		s.log.Errorf("find merchant: %v", err)
		return "", fmt.Errorf("find merchant: %w", err)
	}

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		s.log.Errorf("read file: %v", err)
		return "", fmt.Errorf("read file: %w", err)
	}

	mime := http.DetectContentType(buf)
	if !allowedMIMEs[mime] {
		s.log.Errorf("unsupported file type: %s", mime)
		return "", fmt.Errorf("unsupported file type: %s", mime)
	}

	key, err := s3client.HashKey(file)
	if err != nil {
		s.log.Errorf("hash file: %v", err)
		return "", fmt.Errorf("hash file: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		s.log.Errorf("seek file: %v", err)
		return "", fmt.Errorf("seek file: %w", err)
	}

	key = fmt.Sprintf("%s.%s", key, MIMEExtensions[mime])

	if err := s.s3.Upload(ctx, key, mime, file); err != nil {
		s.log.Errorf("upload: %v", err)
		return "", fmt.Errorf("upload: %w", err)
	}

	url := s.s3.URL(key)
	merchant.PictureURL = &url
	if err := s.merchants.Update(ctx, merchant); err != nil {
		s.log.Errorf("update merchant: %v", err)
		return "", fmt.Errorf("update merchant: %w", err)
	}

	return url, nil
}

func (s *Service) DeletePicture(ctx context.Context, merchantID uuid.UUID) error {
	merchant, err := s.FindMerchantByID(ctx, merchantID)
	if err != nil {
		return err
	}

	if merchant.PictureURL == nil {
		return nil
	}

	key := s.s3.GetKey(*merchant.PictureURL)

	if err := s.s3.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	merchant.PictureURL = nil
	if err := s.merchants.Update(ctx, merchant); err != nil {
		return fmt.Errorf("update merchant: %w", err)
	}

	return nil
}
