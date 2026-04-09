package handlers

import (
	"hello-go-project/myfatoorah/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var MyFatoorahErrors = map[string]string{
	"MF001": "3DS authentication failed (Invalid password, enrollment issue, or bank technical error).",
	"MF002": "The issuer bank has declined the transaction (Invalid card details, insufficient funds, or card not enabled for online use).",
	"MF003": "The transaction has been blocked from the gateway (Unsupported card, fraud detection, or security blocking).",
	"MF004": "Insufficient funds",
	"MF005": "Session timeout",
	"MF006": "Transaction canceled",
	"MF007": "The card is expired.",
	"MF008": "The card issuer doesn't respond.",
	"MF009": "Denied by Risk",
	"MF010": "Wrong Security Code",
	"MF020": "Unspecified Failure",
}

func MyFatoorahWebhookHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.WebhookRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON structure"})
			return
		}

		invoiceID, _ := strconv.Atoi(payload.Data.Invoice.Id)
		invoiceVal, _ := strconv.ParseFloat(payload.Data.Amount.ValueInDisplayCurrency, 64)

		// 1. MetaData & Error Extraction
		var sessionID *uuid.UUID
		var customerID *string
		if payload.Data.Invoice.MetaData != nil {
			meta := payload.Data.Invoice.MetaData
			customerID = meta.UDF2
			if meta.UDF1 != nil && *meta.UDF1 != "" {
				if parsedUUID, err := uuid.Parse(*meta.UDF1); err == nil {
					sessionID = &parsedUUID
				}
			}
		}

		var errCode *string
		if payload.Data.Transaction.Status == "FAILED" && payload.Data.Transaction.Error != nil {
			code := payload.Data.Transaction.Error.Code
			errCode = &code
		}

		// 2. Map to your Entity
		newTransaction := models.Transaction{
			ID:           uuid.New(),
			SessionID:    sessionID,
			InvoiceID:    invoiceID,
			OrderID:      payload.Data.Invoice.Reference,
			Status:       payload.Data.Transaction.Status,
			InvoiceValue: invoiceVal,
			// Use the specific PaymentId from the attempt as your unique reference
			Reference:  payload.Data.Transaction.PaymentId,
			CustomerID: customerID,
			ErrorCode:  errCode,
		}

		// 3. CORRECT Upsert Logic
		// Target 'reference' (which is the MyFatoorah PaymentId) instead of invoice_id.
		// This allows multiple rows for the same InvoiceID (History)
		// but prevents the SAME attempt from being saved twice if the webhook retries.
		err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference"}}, // PaymentId is the unique anchor for the attempt
			DoUpdates: clause.AssignmentColumns([]string{
				"status",
				"error_code",
				"updated_at",
			}),
		}).Create(&newTransaction).Error

		if err != nil {
			log.Printf("DATABASE ERROR: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 4. Business Logic Trigger
		if payload.Data.Transaction.Status == "SUCCESS" {
			log.Printf("Order %s has been PAID successfully", payload.Data.Invoice.Reference)
			// Call your service here to unlock features/complete order
		}

		c.JSON(http.StatusOK, gin.H{"message": "Processed successfully"})
	}
}
