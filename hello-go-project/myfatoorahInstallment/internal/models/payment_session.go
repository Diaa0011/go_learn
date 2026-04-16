package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentSession struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;index" json:"customer_id"`
	InvoiceID  int       `gorm:"type:int;index" json:"invoice_id"`
	Customer   Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`

	// Session B: From Callback (Frontend)
	ExecutionSID string `gorm:"index" json:"execution_sid"`

	Amount float64 `gorm:"type:decimal(10,2)" json:"amount"`
	Status string  `gorm:"type:varchar(20);default:'PENDING'" json:"status"` // PENDING, TOKENIZED, COMPLETED, FAILED

	// Installment Parameters (Stored here until the first payment succeeds)
	TotalIterations int `json:"total_iterations"`
	IntervalDays    int `json:"interval_days"`

	CreatedAt time.Time `json:"created_at"`
}
