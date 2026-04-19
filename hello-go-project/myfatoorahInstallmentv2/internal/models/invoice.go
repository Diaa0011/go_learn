package models

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	PlanID       uuid.UUID     `gorm:"type:uuid;index;not null" json:"plan_id"`
	Amount       float64       `gorm:"not null" json:"amount"`
	DueDate      time.Time     `gorm:"index;not null" json:"due_date"`
	Status       string        `gorm:"default:'pending'" json:"status"`
	MFInvoiceID  *string       `gorm:"index" json:"mf_invoice_id"`
	Transactions []Transaction `gorm:"foreignKey:InvoiceID" json:"transactions,omitempty"`
	PaidAt       *time.Time    `json:"paid_at"`
}
