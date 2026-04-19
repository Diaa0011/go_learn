package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentSession struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID    string    `gorm:"index;not null" json:"customer_id"`
	MFSessionID   string    `gorm:"uniqueIndex;not null" json:"mf_session_id"`
	Status        string    `gorm:"default:'initiated'" json:"status"`
	EncryptionKey string    `gorm:"not null" json:"encryption_key"`
	ExpiresAt     time.Time `gorm:"index" json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}
