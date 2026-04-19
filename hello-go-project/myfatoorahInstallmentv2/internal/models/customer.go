package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string            `gorm:"uniqueIndex;not null" json:"email"`
	Mobile    string            `gorm:"index" json:"mobile"`
	Name      string            `json:"name"`
	Plans     []InstallmentPlan `gorm:"foreignKey:CustomerID" json:"plans,omitempty"`
	Tokens    []PaymentToken    `gorm:"foreignKey:CustomerID" json:"tokens,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}
