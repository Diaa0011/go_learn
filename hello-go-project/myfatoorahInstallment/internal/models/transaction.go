package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Link to the parent installment
	InstallmentID uuid.UUID `gorm:"type:uuid;index" json:"installment_id"`

	// Link to the bridge session (Only for the 1st payment)
	PaymentSessionID *uuid.UUID `gorm:"type:uuid;index" json:"payment_session_id"`

	// Unique payment ID from MyFatoorah for THIS specific hit
	MFPaymentID string `gorm:"uniqueIndex" json:"mf_payment_id"`

	IterationNumber int     `json:"iteration_number"` // Which round is this? (1, 2, 3, or 4)
	Amount          float64 `gorm:"type:decimal(10,2)" json:"amount"`
	Status          string  `gorm:"type:varchar(20)" json:"status"` // SUCCESS, FAILED

	ErrorCode *string `json:"error_code"`

	CreatedAt time.Time `json:"created_at"`
}
