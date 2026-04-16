package handlers

import (
	"bytes"
	"hello-go-project/myfatoorahInstallment/internal/dto/response"
	"hello-go-project/myfatoorahInstallment/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MyFatoorahWebhook handles the asynchronous payment notification from MyFatoorah.
// @Summary      Handle MyFatoorah Webhook
// @Description  Receives payment status updates. If 'Captured', it creates the Installment, Transaction, and updates the Invoice.
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Router       /index.php [post]
// @Security     MyFatoorahAuth
func MyFatoorahWebhook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. READ AND LOG THE RAW RESPONSE
		rawData, err1 := c.GetRawData()
		if err1 != nil {
			log.Printf("[Webhook Error] Could not read raw data: %v", err1)
			c.Status(http.StatusBadRequest)
			return
		}
		log.Printf("[MyFatoorah Webhook Raw Response]: %s", string(rawData))

		// 2. RE-INJECT DATA INTO CONTEXT
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

		// 3. SECURE SIGNATURE CHECK
		// secret := c.GetHeader("MyFatoorah-Signature")
		// key := os.Getenv("WEBHOOK_SECRET")
		// // Corrected: added "!" because if it's NOT valid, return 401
		// if !utils.ValidateMyFatoorahSignature(rawData, key, secret) {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		// 	return
		// }

		// 4. BIND TO UPDATED NESTED STRUCT
		var payload response.WebhookPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("[Webhook Error] JSON Binding: %v", err)
			c.Status(http.StatusBadRequest)
			return
		}

		// 5. EXTRACT VALUES FROM NESTED JSON
		// Convert String IDs from JSON to Integers for your Database
		mfInvoiceID, _ := strconv.Atoi(payload.Data.Invoice.Id)
		mfPaymentID := payload.Data.Transaction.PaymentId
		txStatus := payload.Data.Transaction.Status
		invoiceValue, _ := strconv.ParseFloat(payload.Amount.ValueInDisplayCurrency, 64)

		sessionUUID := payload.Data.Invoice.ExternalIdentifier

		// This is the UUID you passed in 'CustomerReference' or 'ExternalIdentifier'

		// 6. Only proceed if the payment was successful
		if txStatus != "SUCCESS" && payload.Data.Invoice.Status != "PAID" {
			log.Printf("Payment not successful: Status %s for Invoice %d", txStatus, mfInvoiceID)
			c.Status(http.StatusOK)
			return
		}

		// 7. Find the Session using the UUID (ExternalIdentifier)
		// This is safer than looking up by Invoice ID
		var session models.PaymentSession
		if err := db.Where("id = ?", sessionUUID).First(&session).Error; err != nil {
			log.Printf("Session not found for UUID: %s (MF Invoice: %d)", sessionUUID, mfInvoiceID)
			c.Status(http.StatusOK)
			return
		}

		// 8. ATOMIC EXECUTION
		err := db.Transaction(func(tx *gorm.DB) error {

			// A. Create the Installment Plan
			installment := models.Installment{
				ID:          uuid.New(),
				CustomerID:  session.CustomerID,
				TotalAmount: invoiceValue,
				Status:      "ACTIVE",

				// Note: If RecurringId isn't in this webhook,
				// you might need to fetch it from a 'GetPaymentStatus' call
				// or check payload.Data.Transaction.Card.Token
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&installment).Error; err != nil {
				return err
			}

			// B. Create the Transaction
			newTx := models.Transaction{
				ID:               uuid.New(),
				InstallmentID:    installment.ID,
				PaymentSessionID: &session.ID,
				MFInvoiceID:      mfInvoiceID,
				MFPaymentID:      mfPaymentID,
				IterationNumber:  1,
				Amount:           invoiceValue,
				Status:           "SUCCESS",
				CreatedAt:        time.Now(),
			}
			if err := tx.Create(&newTx).Error; err != nil {
				return err
			}

			// C. Update the Invoice
			if err := tx.Model(&models.Invoice{}).
				Where("mf_invoice_id = ?", mfInvoiceID).
				Updates(map[string]interface{}{
					"status":         "PAID",
					"installment_id": installment.ID,
					"paid_at":        time.Now(),
				}).Error; err != nil {
				return err
			}

			// D. Close the Session
			session.Status = "COMPLETED"
			if err := tx.Save(&session).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("Webhook DB Transaction Failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal processing error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Infrastructure created successfully"})
	}
}
