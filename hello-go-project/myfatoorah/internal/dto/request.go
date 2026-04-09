package dto

import "time"

type CreateSessionRequest struct {
	Amount         float64    `json:"amount" binding:"required"`
	OrderID        *string    `json:"order_id"`    // Nullable
	CustomerID     *string    `json:"customer_id"` // Nullable
	RedirectionUrl string     `json:"redirection_url"`
	SessionExpiry  *time.Time `json:"session_expiry"` // Nullable UTC time
}
