package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	// Internal Primary Key
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// MyFatoorah Data
	MyFatoorahSessionID string    `gorm:"uniqueIndex;not null" json:"session_id"`
	SessionExpiry       time.Time `json:"session_expiry"`
	EncryptionKey       string    `json:"encryption_key"`
	OperationType       string    `json:"operation_type"` // e.g., "PAY"

	// Order Details
	Amount   float64 `gorm:"type:decimal(10,2)" json:"amount"`
	Currency string  `gorm:"type:varchar(5)" json:"currency"`

	// Customer Details (Mapped from the Customer object in response)
	CustomerName      string `json:"customer_name"`
	CustomerReference string `json:"customer_reference"` // This is your internal "456123789"
	CustomerEmail     string `json:"customer_email"`

	// Relations
	Transactions []Transaction `gorm:"foreignKey:SessionID;references:ID" json:"transactions"`

	// Metadata
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID    *uuid.UUID     `gorm:"type:uuid;index" json:"session_id"` // Change string to sql.NullString	InvoiceID    int       `gorm:"uniqueIndex" json:"invoice_id"`
	InvoiceID    int            `gorm:"index" json:"invoice_id"`
	Reference    string         `json:"reference"` // myfatoorah reference id
	OrderID      string         `gorm:"index" json:"order_id"`
	MerchantID   *string        `json:"merchant_id"`
	CustomerID   *string        `json:"customer_id"`
	Status       string         `gorm:"type:varchar(20)" json:"status"`
	InvoiceValue float64        `gorm:"type:decimal(10,2)" json:"invoice_value"`
	ErrorCode    *string        `gorm:"type:varchar(20)" json:"-"`
	ErrorMessage string         `gorm:"-" json:"error_message"` // This IS NOT in the DB
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// WebhookRequest is the top-level structure sent by MyFatoorah
type WebhookRequest struct {
	Event Event       `json:"Event"` // e.g., "TransactionsStatusChanged"
	Data  WebhookData `json:"Data"`  // The actual transaction details
}

type Event struct {
	Name string `json:"Name"`
}

// WebhookData contains the specific payment information
type WebhookData struct {
	Invoice struct {
		Id               string `json:"Id"`
		Reference        string `json:"Reference"`
		UserDefinedField string `json:"UserDefinedField"`
		MetaData         *struct {
			UDF1 *string `json:"UDF1"`
			UDF2 *string `json:"UDF2"`
			UDF3 *string `json:"UDF3"`
		} `json:"MetaData"`
	} `json:"Invoice"`
	Transaction struct {
		Status    string `json:"Status"`
		PaymentId string `json:"PaymentId"`
		// Add this part to capture the error details
		Error *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"Transaction"`
	Customer struct {
		Email string `json:"Email"`
	} `json:"Customer"`
	Amount struct {
		ValueInDisplayCurrency string `json:"ValueInDisplayCurrency"`
	} `json:"Amount"`
}

// MyFatoorahSessionResponse maps the actual response from MyFatoorah API
type MyFatoorahSessionResponse struct {
	IsSuccess bool   `json:"IsSuccess"`
	Message   string `json:"Message"`
	Data      struct {
		SessionId     string `json:"SessionId"`
		SessionExpiry string `json:"SessionExpiry"`
		EncryptionKey string `json:"EncryptionKey"`
	} `json:"Data"`
}
