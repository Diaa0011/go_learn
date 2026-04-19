package models

import (
	"time"

	"github.com/google/uuid"
)

type InstallmentPlan struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID         uuid.UUID  `gorm:"type:uuid;index;not null" json:"customer_id"`
	TotalAmount        float64    `gorm:"not null" json:"total_amount"`
	RemainingAmount    float64    `gorm:"not null" json:"remaining_amount"`
	Status             string     `gorm:"default:'pending'" json:"status"`
	NextPayTime        *time.Time `gorm:"index" json:"next_pay_time"`
	LastPayTime        *time.Time `json:"last_pay_time"`
	ExternalIdentifier string     `gorm:"uniqueIndex" json:"external_identifier"`
	Invoices           []Invoice  `gorm:"foreignKey:PlanID" json:"invoices,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
