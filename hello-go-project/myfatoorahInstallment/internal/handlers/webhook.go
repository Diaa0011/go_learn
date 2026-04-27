package handlers

import (
	"fmt"
	"hello-go-project/myfatoorahInstallment/internal/models"
	"hello-go-project/myfatoorahInstallment/internal/utils"
	"log"
	"net/http"
	"os"
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

		var rawPayload map[string]interface{}
		if err := c.ShouldBindJSON(&rawPayload); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		myFatoorahPayload := utils.BuildSignture(rawPayload)
		log.Printf("Payload for Signature: %s", myFatoorahPayload)
		bytedPayload := []byte(myFatoorahPayload)

		signature := c.GetHeader("MyFatoorah-Signature")
		secret := os.Getenv("WEBHOOK_SECRET")
		if !utils.ValidateMyFatoorahSignature(bytedPayload, secret, signature) {
			log.Printf("body: %s", rawPayload)
			log.Printf("Invalid signature for webhook: %s", signature)
			c.Status(http.StatusUnauthorized)
			return
		}
		log.Printf("Webhook signature validated successfully")

		// 1. Identify the Event Type
		// Code 1 is in a nested "Event" object, Code 5 is often top-level or in "EventCode"
		eventCode := extractEventCode(rawPayload)

		switch eventCode {
		case 1: // PAYMENT_STATUS_CHANGED (The First/Manual Payment)
			handleInitialPayment(c, db, rawPayload)
		case 5: // RECURRING_UPDATES (Automated Installments)
			handleRecurringUpdate(c, db, rawPayload)
		default:
			log.Printf("Unhandled Event Code: %v", eventCode)
			c.Status(http.StatusOK)
		}
	}
}

// Helper to handle the FIRST payment (Event Code 1)
func handleInitialPayment(c *gin.Context, db *gorm.DB, raw map[string]interface{}) {
	data := raw["Data"].(map[string]interface{})
	invoice := data["Invoice"].(map[string]interface{})
	transaction := data["Transaction"].(map[string]interface{})

	// SAFE CONVERSION:
	var mfInvoiceID int
	if idStr, ok := invoice["Id"].(string); ok {
		mfInvoiceID, _ = strconv.Atoi(idStr)
	} else if idNum, ok := invoice["Id"].(float64); ok {
		mfInvoiceID = int(idNum)
	}

	mfPaymentID := transaction["PaymentId"].(string)

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Find the Invoice Draft we created in Phase A
		var inv models.Invoice
		if err := tx.Where("mf_invoice_id = ?", mfInvoiceID).First(&inv).Error; err != nil {
			return fmt.Errorf("invoice draft not found: %d", mfInvoiceID)
		}

		// 2. Find the associated Installment (linked via foreign key or logic)
		var inst models.Installment
		if err := tx.Where("id = ?", inv.InstallmentID).First(&inst).Error; err != nil {
			// Alternatively, if you didn't link them yet:
			// tx.Where("customer_id = ? AND status = 'PENDING'", inv.CustomerID).First(&inst)
			return err
		}

		// 3. Activate records
		tx.Model(&inv).Updates(map[string]interface{}{
			"status":  "PAID",
			"paid_at": time.Now(),
		})

		tx.Model(&inst).Updates(map[string]interface{}{
			"status":            "ACTIVE",
			"current_iteration": 1,
		})

		// 4. Create first history entry
		tx.Create(&models.Transaction{
			ID:              uuid.New(),
			InstallmentID:   inst.ID,
			MFInvoiceID:     mfInvoiceID,
			MFPaymentID:     mfPaymentID,
			IterationNumber: 1,
			Amount:          inv.Amount,
			Status:          "SUCCESS",
		})

		// 5. Update Session for UI feedback
		tx.Model(&models.PaymentSession{}).Where("invoice_id = ?", mfInvoiceID).Update("status", "COMPLETED")

		return nil
	})

	if err != nil {
		log.Printf("[Webhook Code 1] Error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// Helper to handle AUTOMATED installments (Event Code 5)
func handleRecurringUpdate(c *gin.Context, db *gorm.DB, raw map[string]interface{}) {
	data := raw["Data"].(map[string]interface{})
	recurring := data["Recurring"].(map[string]interface{})
	payment := data["Payment"].(map[string]interface{})

	mfRecurringID := recurring["Id"].(string)
	mfInvoiceID := int(payment["Invoice"].(map[string]interface{})["Id"].(float64))
	nextPayDate := recurring["NextPayDate"].(string)

	err := db.Transaction(func(tx *gorm.DB) error {
		var inst models.Installment
		if err := tx.Where("mf_recurring_id = ?", mfRecurringID).First(&inst).Error; err != nil {
			return err
		}

		// Prevent Duplicate Processing
		var count int64
		tx.Model(&models.Transaction{}).Where("mf_invoice_id = ?", mfInvoiceID).Count(&count)
		if count > 0 {
			return nil
		}

		// Update Installment Progress
		inst.CurrentIteration += 1
		if t, err := time.Parse(time.RFC3339, nextPayDate); err == nil {
			inst.NextBillingDate = t
		}
		tx.Save(&inst)

		// Create the new Transaction record for this month
		tx.Create(&models.Transaction{
			ID:              uuid.New(),
			InstallmentID:   inst.ID,
			MFInvoiceID:     mfInvoiceID,
			Amount:          inst.IterationAmount,
			Status:          "SUCCESS",
			IterationNumber: inst.CurrentIteration,
		})

		return nil
	})

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func extractEventCode(raw map[string]interface{}) int {
	// Check if it's the structure of Event 1
	if event, ok := raw["Event"].(map[string]interface{}); ok {
		return int(event["Code"].(float64))
	}
	// Check if it's the structure of Event 5
	if code, ok := raw["EventCode"].(float64); ok {
		return int(code)
	}
	return 0
}
