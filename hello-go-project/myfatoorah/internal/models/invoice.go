package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentSource string

const (
	SourceHosted   PaymentSource = "HOSTED"
	SourceEmbedded PaymentSource = "EMBEDDED"
)

type PaymentType string

const (
	TypeOneTime     PaymentType = "ONE_TIME"
	TypeInstallment PaymentType = "INSTALLMENT"
)

type Invoice struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID    string    `gorm:"uniqueIndex;not null" json:"order_id"` // Internal ID
	TotalValue float64   `gorm:"type:decimal(10,2)" json:"total_value"`
	Currency   string    `gorm:"type:varchar(5);default:'KWD'" json:"currency"`
	Status     string    `gorm:"type:varchar(20);default:'PENDING'" json:"status"` // PENDING, PAID, EXPIRED

	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`

	Source PaymentSource `gorm:"type:varchar(20);not null;default:'HOSTED'" json:"source"`
	Type   PaymentType   `gorm:"type:varchar(20);not null;default:'ONE_TIME'" json:"type"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:InvoiceID" json:"transactions"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
