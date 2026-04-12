package handlers

import (
	"bytes"
	"encoding/json"
	"hello-go-project/myfatoorah/internal/dto"
	"hello-go-project/myfatoorah/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateInvoiceHandler godoc
// @Summary      Initiate Hosted Payment (Invoice/Link)
// @Description  Creates an internal invoice and generates a MyFatoorah payment link.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreatePaymentRequest  true  "Invoice Initiation Details"
// @Success      201      {object}  dto.CreatePaymentResponse
// @Router       /invoices [post]
func CreateInvoiceHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. SAFE EXTRACTION: Prevent nil pointer panic for OrderID
		orderRef := "ORD-" + uuid.New().String()[:8]
		if req.OrderID != nil && *req.OrderID != "" {
			orderRef = *req.OrderID
		}

		// 2. CREATE INTERNAL INVOICE (Parent)
		// We create this before the API call so we have a UUID to send in Metadata
		internalInvoice := models.Invoice{
			ID:            uuid.New(),
			OrderID:       orderRef,
			TotalValue:    req.Amount,
			Currency:      "KWD",
			Status:        "PENDING",
			Source:        models.SourceEmbedded,
			Type:          models.TypeOneTime,
			CustomerName:  "Diaa Dawood", // Your default logic
			CustomerEmail: "diaadawood.mas@gmail.com",
			CreatedAt:     time.Now(),
		}

		if err := db.Create(&internalInvoice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create internal record"})
			return
		}

		// 3. PREPARE MYFATOORAH PAYLOAD
		// Keeping all original values as requested, but adding Metadata bridge
		mfPayload := map[string]interface{}{
			"Order": map[string]interface{}{
				"Amount":             req.Amount,
				"Currency":           "KWD",
				"ExternalIdentifier": orderRef,
			},
			"Customer": map[string]interface{}{
				"Name":  internalInvoice.CustomerName,
				"Email": internalInvoice.CustomerEmail,
				"Mobile": map[string]interface{}{
					"CountryCode": "20",
					"Number":      "101543346",
				},
			},
			"MetaData": map[string]interface{}{
				"UDF1": internalInvoice.ID.String(), // BRIDGE: Links Hosted ID to our UUID
			},
			"NotificationOption": "ALL",
			"IntegrationUrls": map[string]string{
				"Redirection": "https://google.com/payment-callback",
			},
		}

		// 4. CALL MYFATOORAH API
		resp, err := callMyFatoorahAPI(mfPayload)
		if err != nil || resp["IsSuccess"] == false {
			msg := "Gateway error"
			if resp != nil && resp["Message"] != nil {
				msg = resp["Message"].(string)
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": msg})
			return
		}

		data := resp["Data"].(map[string]interface{})

		// 5. SAFE EXTRACTION OF REMOTE IDs
		var mfInvoiceID int
		if val, ok := data["InvoiceId"]; ok {
			switch v := val.(type) {
			case float64:
				mfInvoiceID = int(v)
			case string:
				id, _ := strconv.Atoi(v)
				mfInvoiceID = id
			}
		}

		paymentURL, _ := data["PaymentURL"].(string)

		// // 6. CREATE INITIAL TRANSACTION (Optional but recommended for hosted)
		// // We create one pending transaction to track the MyFatoorah Invoice ID immediately.
		// transaction := models.Transaction{
		// 	ID:           uuid.New(),
		// 	InvoiceID:    internalInvoice.ID,
		// 	MFInvoiceID:  mfInvoiceID,
		// 	OrderID:      orderRef,
		// 	Status:       "PENDING",
		// 	InvoiceValue: req.Amount,
		// 	CreatedAt:    time.Now(),
		// }
		// db.Create(&transaction)

		// 7. RETURN RESPONSE
		c.JSON(http.StatusCreated, gin.H{
			"message":             "Payment link generated",
			"internal_invoice_id": internalInvoice.ID,
			"mf_invoice_id":       mfInvoiceID,
			"payment_url":         paymentURL,
		})
	}
}

func callMyFatoorahAPI(payload interface{}) (map[string]interface{}, error) {
	apiURL := "https://apitest.myfatoorah.com/v3/payments"
	token := "SK_KWT_Uw42KMHRKdVNZLHC4KEBTkHxOMfVsQi4txBkvQ1A1iAbVDdfPRhFsMQLVLkLXY3S"

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetInvoice godoc
// @Summary      Get a specific invoice
// @Description  Retrieves an invoice by its UUID, including all associated payment sessions and their transaction attempts.
// @Tags         invoices
// @Produce      json
// @Param        id   path      string  true  "Invoice UUID"
// @Success      200  {object}  models.Invoice
// @Failure      404  {object}  map[string]string "error: Invoice not found"
// @Router       /invoices/{id} [get]
func GetInvoice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var invoice models.Invoice

		// Preload Transactions directly linked to this Invoice UUID
		if err := db.Preload("Transactions").Where("id = ?", id).First(&invoice).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		c.JSON(http.StatusOK, invoice)
	}
}

// GetAllInvoices godoc
// @Summary      List all invoices
// @Description  Retrieves all invoices with their nested session and transaction history.
// @Tags         invoices
// @Produce      json
// @Success      200  {array}   models.Invoice
// @Failure      500  {object}  map[string]string "error: database error"
// @Router       /invoices [get]
func GetAllInvoices(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var invoices []models.Invoice
		// Preload all transaction attempts for every invoice
		if err := db.Preload("Transactions").Find(&invoices).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, invoices)
	}
}
