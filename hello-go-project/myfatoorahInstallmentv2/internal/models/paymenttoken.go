package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;index;not null" json:"customer_id"`
	Token      string    `gorm:"not null" json:"token"`
	CardBrand  string    `json:"card_brand"`
	CardLast4  string    `json:"card_last4"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}
