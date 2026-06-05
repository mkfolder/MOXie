package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58/base58"
	"gorm.io/gorm"
)

type Order struct {
	BaseModel
	MerchantID         uuid.UUID       `json:"merchant_id" gorm:"not null"`
	Merchant           Merchant        `json:"-" gorm:"foreignKey:MerchantID;constraint:OnDelete:CASCADE"`
	Address            string          `json:"address" gorm:"-"`
	Memo               string          `json:"memo" gorm:"not null"`
	RawRequestedAmount uint64          `json:"raw_requested_amount"`
	RawPaidAmount      *uint64         `json:"raw_paid_amount" gorm:"default:0"`
	TxHash             *string         `json:"tx_hash" gorm:"default:null"`
	PaidAt             *time.Time      `json:"paid_at" gorm:"default:null"`
	CustomData         json.RawMessage `json:"custom_data" gorm:"type:jsonb;not null"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.Address == "" {
		return errors.New("invalid order receiver address")
	}
	if o.MerchantID == uuid.Nil {
		return errors.New("invalid order merchant id")
	}
	return nil
}

func (o *Order) AfterCreate(tx *gorm.DB) error {
	o.Memo = base58.Encode(o.ID[:])
	return tx.Save(o).Error
}
