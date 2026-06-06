package models

import (
	"errors"

	"github.com/mkfolder/moxie/internal/common"
	"github.com/mr-tron/base58/base58"
	"gorm.io/gorm"
)

type Merchant struct {
	BaseModel
	APIKey      string  `json:"api_key" gorm:"unique;not null"`
	Email       string  `json:"email" gorm:"unique;not null"`
	Address     string  `json:"address" gorm:"not null"`
	PasswdHash  []byte  `json:"-" gorm:"not null"`
	WebhookURL  *string `json:"webhook_url" gorm:"default:null"`
	TOTPSecret  string  `json:"-" gorm:"default:''"`
	TOTPEnabled bool    `json:"totp_enabled" gorm:"default:false"`
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.Email == "" {
		return errors.New("invalid merchant email")
	}
	if err := common.ValidateEmail(m.Email); err != nil {
		return err
	}
	if len(m.PasswdHash) == 0 {
		return errors.New("invalid merchant password hash")
	}
	return nil
}

func (m *Merchant) AfterCreate(tx *gorm.DB) error {
	m.APIKey = base58.Encode(m.ID[:])
	return tx.Save(m).Error
}
