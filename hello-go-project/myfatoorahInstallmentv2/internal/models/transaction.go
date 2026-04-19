package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	InvoiceID       uuid.UUID `gorm:"index;not null"`
	MFTrackID       string    `gorm:"index"` // The TrackId from MyFatoorah
	GatewayStatus   string    // Success, Captured, Failed, Canceled
	Amount          float64
	TransactionDate time.Time
	ErrorCode       string
	RawResponse     string `gorm:"type:text"` // Store the full JSON for debugging
}
