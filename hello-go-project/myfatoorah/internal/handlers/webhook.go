package handlers

import (
	"hello-go-project/myfatoorah/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

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

		mfInvoiceID, _ := strconv.Atoi(payload.Data.Invoice.Id)
		invoiceVal, _ := strconv.ParseFloat(payload.Data.Amount.ValueInDisplayCurrency, 64)

		// 1. Extract internal Invoice UUID from Metadata (UDF1)
		var internalInvoiceID *uuid.UUID
		if payload.Data.Invoice.MetaData != nil && payload.Data.Invoice.MetaData.UDF1 != nil {
			parsedID, err := uuid.Parse(*payload.Data.Invoice.MetaData.UDF1)
			if err != nil {
				log.Printf("WEBHOOK ERROR: Could not parse UDF1 UUID: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Metadata ID"})
				return
			}
			internalInvoiceID = &parsedID
		}

		// 2. Find the Session linked to this Invoice
		// This ensures the Transaction gets a valid SessionID for the relationship
		var sessionID *uuid.UUID
		var session models.Session
		if err := db.Where("invoice_id = ?", internalInvoiceID).First(&session).Error; err == nil {
			sessionID = &session.ID
		} else {
			// log.Printf("WEBHOOK WARNING: No session found for Invoice %s. Proceeding without SessionID.", internalInvoiceID)
		}

		// 3. Extract Error Code if attempt failed
		var errCode *string
		if payload.Data.Transaction.Status == "FAILED" && payload.Data.Transaction.Error != nil {
			code := payload.Data.Transaction.Error.Code
			errCode = &code
		}

		// 4. Prepare the Transaction record
		newTransaction := models.Transaction{
			ID:           uuid.New(),
			InvoiceID:    internalInvoiceID,
			SessionID:    sessionID, // This is now correctly linked!
			MFInvoiceID:  mfInvoiceID,
			Reference:    payload.Data.Transaction.PaymentId,
			OrderID:      payload.Data.Invoice.Reference,
			Status:       payload.Data.Transaction.Status,
			InvoiceValue: invoiceVal,
			ErrorCode:    errCode,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// 5. Execute DB updates in a transaction
		err := db.Transaction(func(tx *gorm.DB) error {
			// Upsert: Create if new attempt, Update if retry
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "reference"}}, // PaymentId uniqueness
				DoUpdates: clause.AssignmentColumns([]string{"status", "error_code", "updated_at"}),
			}).Create(&newTransaction).Error; err != nil {
				return err
			}

			// Update parent Invoice status if payment was successful
			if payload.Data.Transaction.Status == "SUCCESS" {
				if err := tx.Model(&models.Invoice{}).
					Where("id = ?", internalInvoiceID).
					Update("status", "PAID").Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("WEBHOOK DB ERROR: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Processed successfully"})
	}
}
