package models

import (
	"time"

	"github.com/google/uuid"
)

type Installment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;index" json:"customer_id"`

	// The "Master ID" from MyFatoorah recurring response
	MFRecurringID string `gorm:"uniqueIndex" json:"mf_recurring_id"`

	Status          string  `gorm:"type:varchar(20)" json:"status"` // ACTIVE, COMPLETED, CANCELED, PAST_DUE
	TotalAmount     float64 `gorm:"type:decimal(10,2)" json:"total_amount"`
	IterationAmount float64 `gorm:"type:decimal(10,2)" json:"iteration_amount"`

	TotalIterations  int `json:"total_iterations"`
	CurrentIteration int `json:"current_iteration"` // e.g., 1, 2, 3...

	NextBillingDate time.Time `json:"next_billing_date"`

	// The card token for future hits (if needed for manual retries)
	CardToken string `json:"card_token"`

	Transactions []Transaction `gorm:"foreignKey:InstallmentID" json:"transactions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
