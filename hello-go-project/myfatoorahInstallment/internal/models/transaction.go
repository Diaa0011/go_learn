package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	InstallmentID    uuid.UUID  `gorm:"type:uuid;index" json:"installment_id"`
	PaymentSessionID *uuid.UUID `gorm:"type:uuid;index" json:"payment_session_id"`

	// Add this for Webhook lookup!
	MFInvoiceID int `gorm:"uniqueIndex" json:"mf_invoice_id"`

	MFPaymentID     string    `gorm:"uniqueIndex" json:"mf_payment_id"`
	IterationNumber int       `json:"iteration_number"`
	Amount          float64   `gorm:"type:decimal(10,2)" json:"amount"`
	Status          string    `gorm:"type:varchar(20)" json:"status"`
	ErrorCode       *string   `json:"error_code"`
	CreatedAt       time.Time `json:"created_at"`
}
