package handlers

import (
	"bytes"
	"encoding/json"
	"hello-go-project/myfatoorah/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 1. Define the struct to hold your DB connection
type PaymentHandler struct {
	DB *gorm.DB
}

// 2. Make CreatePayment a method of the struct
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var input struct {
		Amount       float64 `json:"amount" binding:"required"`
		CustomerName string  `json:"customer_name" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		OrderID      string  `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mfPayload := map[string]interface{}{
		"Order": map[string]interface{}{
			"Amount":             input.Amount,
			"Currency":           "KWD",
			"ExternalIdentifier": input.OrderID,
		},
		"Customer": map[string]interface{}{
			"Name":  input.CustomerName,
			"Email": input.Email,
		},
		"NotificationOption": "ALL",
		"IntegrationUrls": map[string]string{
			"Redirection": "https://google.com/payment-callback",
		},
	}

	// 3. Call the helper (now a method or stays internal)
	resp, err := h.callMyFatoorah(mfPayload)
	if err != nil || resp["IsSuccess"] == false {
		message := "Unknown error"
		if resp != nil && resp["Message"] != nil {
			message = resp["Message"].(string)
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "Payment gateway error", "details": message})
		return
	}

	data := resp["Data"].(map[string]interface{})

	// Safe type assertions
	invoiceID := int(data["InvoiceId"].(float64))
	paymentURL := data["PaymentURL"].(string)

	// Handle InvoiceReference safely as it might be null/empty in some v3 responses
	reference := ""
	if val, ok := data["InvoiceReference"].(string); ok {
		reference = val
	}

	// 4. h.DB is now accessible because this is a method of PaymentHandler
	transaction := models.Transaction{
		ID:           uuid.New(),
		InvoiceID:    invoiceID,
		OrderID:      input.OrderID,
		Status:       "PENDING",
		InvoiceValue: input.Amount,
		Reference:    reference,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Payment link generated",
		"payment_url": paymentURL,
		"invoice_id":  invoiceID,
	})
}

// 5. Helper function as a method to keep the token/URL logic encapsulated
func (h *PaymentHandler) callMyFatoorah(payload interface{}) (map[string]interface{}, error) {
	apiURL := "https://apitest.myfatoorah.com/v3/payments"
	token := "SK_KWT_Uw42KMHRKdVNZLHC4KEBTkHxOMfVsQi4txBkvQ1A1iAbVDdfPRhFsMQLVLkLXY3S"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}


// internals/handlers/transactions.go
package handlers

import (
	"hello-go-project/myfatoorah/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllTransactions godoc
// @Summary      Get All Transactions
// @Description  Retrieves a list of all payment transactions with search, filter, and sorting.
// @Tags         transactions
// @Produce      json
// @Param        search    query     string  false  "Search by Order ID or Reference"
// @Param        status    query     string  false  "Filter by Status (e.g., SUCCESS, FAILED)"
// @Param        sort_by   query     string  false  "Field to sort by (default: created_at)"
// @Param        order     query     string  false  "Order direction (asc, desc)"
// @Success      200 {array}   models.Transaction
// @Failure      500 {object}  map[string]string
// @Router       /transactions [get]
func GetAllTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactions []models.Transaction
		query := db.Model(&models.Transaction{})

		// --- 1. SEARCH ---
		// Searches across multiple fields if 'search' query param is provided
		if search := c.Query("search"); search != "" {
			searchTerm := "%" + search + "%"
			query = query.Where(
				"order_id ILIKE ? OR status ILIKE ? OR reference ILIKE ?",
				searchTerm, searchTerm, searchTerm,
			)
		}

		// --- 2. FILTERING ---
		// Exact match filters
		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", status)
		}
		if sessionID := c.Query("session_id"); sessionID != "" {
			query = query.Where("session_id = ?", sessionID)
		}

		// --- 3. SORTING ---
		sortBy := c.DefaultQuery("sort_by", "created_at") // Default sort field
		orderDir := c.DefaultQuery("order", "desc")       // Default order direction

		// Validate direction to prevent SQL injection
		if orderDir != "asc" && orderDir != "desc" {
			orderDir = "desc"
		}

		// Build the order string, e.g., "created_at desc"
		query = query.Order(sortBy + " " + orderDir)

		// --- 4. EXECUTION ---
		if err := query.Find(&transactions).Error; err != nil {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// --- 5. HYDRATE VIRTUAL FIELD ---
		for i := range transactions {
			if transactions[i].ErrorCode != nil {
				if msg, exists := MyFatoorahErrors[*transactions[i].ErrorCode]; exists {
					transactions[i].ErrorMessage = msg
				}
			}
		}

		c.JSON(http.StatusOK, transactions)
	}
}
