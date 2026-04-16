package models

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID    uuid.UUID `gorm:"type:uuid;index" json:"customer_id"`
	InstallmentID uuid.UUID `gorm:"type:uuid;index" json:"installment_id"`

	// The reference for the user (e.g., INV-2024-001)
	InvoiceNumber string `gorm:"uniqueIndex" json:"invoice_number"`

	// Links to MyFatoorah
	MFInvoiceID int `gorm:"index" json:"mf_invoice_id"`

	Amount      float64 `gorm:"type:decimal(10,2)" json:"amount"`
	Tax         float64 `gorm:"type:decimal(10,2)" json:"tax"`
	TotalAmount float64 `gorm:"type:decimal(10,2)" json:"total_amount"`

	Status  string     `json:"status"` // DRAFT, SENT, PAID, VOID
	DueDate time.Time  `json:"due_date"`
	PaidAt  *time.Time `json:"paid_at"`

	CreatedAt time.Time `json:"created_at"`
}
