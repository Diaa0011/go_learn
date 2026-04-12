package dto

import "time"

type CreateSessionRequest struct {
	Amount         float64    `json:"amount" binding:"required"`
	OrderID        *string    `json:"order_id"`    // Nullable
	CustomerID     *string    `json:"customer_id"` // Nullable
	RedirectionUrl string     `json:"redirection_url"`
	SessionExpiry  *time.Time `json:"session_expiry"` // Nullable UTC time
}

type CreatePaymentRequest struct {
	Amount     float64 `json:"amount" binding:"required,gt=0" example:"100.00"`
	OrderID    *string `json:"order_id" example:"ORD-2026-001"`
	ExpiryDate *string `json:"expiry_date"`

	// Optional fields with defaults
	CustomerName string `json:"customer_name" example:"Diaa Dawood"`
	Email        string `json:"email" example:"diaadawood.mas@gmail.com"`
	CountryCode  string `json:"country_code" example:"20"`
	MobileNumber string `json:"mobile_number" example:"1015433746"`
}

type CreatePaymentResponse struct {
	Message    string `json:"message" example:"Payment link generated successfully"`
	PaymentURL string `json:"payment_url" example:"https://demo.myfatoorah.com/..."`
	InvoiceID  int    `json:"invoice_id" example:"6646004"`
}
