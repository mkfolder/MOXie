package models

import (
	"errors"
	"strings"

	"github.com/mkfolder/moxie/internal/common"
	"gorm.io/gorm"
)

type Merchant struct {
	BaseModel
	Email            string  `json:"email" gorm:"unique;not null"`
	Username         string  `json:"username" gorm:"not null"`
	Address          *string `json:"address" gorm:"default:null"`
	AvatarURL        *string `json:"avatar_url" gorm:"default:null"`
	IsServiceEnabled bool    `json:"is_service_enabled" gorm:"-"`
	WebhookURL       *string `json:"webhook_url" gorm:"default:null"`
	PasswdHash       []byte  `json:"-" gorm:"not null"`
	HeliusAPIKey     *string `json:"-" gorm:"default:null"`
	TOTPSecret       string  `json:"-" gorm:"default:''"`
	TOTPEnabled      bool    `json:"totp_enabled" gorm:"default:false"`
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.Email == "" {
		return errors.New("invalid merchant email")
	}
	if common.ValidateEmail(m.Email) != nil {
		return errors.New("invalid merchant email")
	}
	if m.Username == "" {
		m.Username = strings.Split(m.Email, "@")[0]
	}
	if len(m.PasswdHash) == 0 {
		return errors.New("invalid merchant password hash")
	}
	return nil
}
