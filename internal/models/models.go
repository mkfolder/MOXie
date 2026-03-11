package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Order struct {
	BaseModel
	Address    string          `json:"address"`
	Memo       string          `json:"-"`
	TxHash     *string         `json:"tx_hash" gorm:"default:null"`
	RawAmount  uint64          `json:"raw_amount"`
	CustomData json.RawMessage `json:"custom_data" gorm:"type:jsonb;not null"`
}
