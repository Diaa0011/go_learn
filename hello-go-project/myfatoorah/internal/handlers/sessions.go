package handlers

import (
	"bytes"
	"encoding/json"
	"hello-go-project/myfatoorah/internal/dto"
	"hello-go-project/myfatoorah/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetAllSessions godoc
// @Summary      Get All Sessions with Transactions
// @Tags         sessions
// @Success      200 {array} models.Session
// @Router       /sessions [get]
func GetAllSessions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessions []models.Session

		// .Preload("Transactions") is the magic that joins the two tables
		if err := db.Preload("Transactions").Find(&sessions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}

// CreateSessionHandler godoc
// @Summary      Create a new payment session
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateSessionRequest  true  "Session Details"
// @Router       /sessions [post]
func CreateSessionHandlerv1(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiURL := "https://apitest.myfatoorah.com/v3/sessions"
		apiKey := "SK_KWT_Uw42KMHRKdVNZLHC4KEBTkHxOMfVsQi4txBkvQ1A1iAbVDdfPRhFsMQLVLkLXY3S"

		sessionId := uuid.New()
		// --- Default Logic ---
		customerName := "Diaa Dawood"
		customerEmail := "diaadawood.mas@gmail.com"
		customerMobile := "01012345678"

		// Use Redirection URL from request or a default fallback
		redirect := req.RedirectionUrl
		if redirect == "" {
			redirect = "https://www.google.com/"
		}

		// Prepare Metadata mapping
		// UDF1: CustomerID, UDF2: CustomerID, UDF3: Placeholder for SessionID
		udf1Value := ""
		if req.CustomerID != nil {
			udf1Value = *req.CustomerID
		}

		payload := map[string]interface{}{
			"PaymentMode": "COMPLETE_PAYMENT",
			"Order": map[string]interface{}{
				"Amount":             req.Amount,
				"Currency":           "KWD",
				"ExternalIdentifier": req.OrderID, // Will be null if OID is nil
			},
			"Customer": map[string]interface{}{
				"Name":      customerName,
				"Email":     customerEmail,
				"Reference": req.OrderID, // Reference ID = Order ID per request
				"Mobile": map[string]string{
					"Number":      customerMobile,
					"CountryCode": "20",
				},
			},
			"IntegrationUrls": map[string]string{
				"Redirection": redirect,
			},
			"MetaData": map[string]interface{}{
				"UDF1": sessionId,
			},
			"SupportedNetworks":       []string{"visa", "masterCard"},
			"SupportedPaymentMethods": []string{"card", "knet", "googlepay", "applepay"},
		}

		jsonData, _ := json.Marshal(payload)
		httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API Call Failed"})
			return
		}
		defer resp.Body.Close()

		var mfResp models.MyFatoorahSessionResponse
		if err := json.NewDecoder(resp.Body).Decode(&mfResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		if !mfResp.IsSuccess {
			c.JSON(http.StatusBadRequest, gin.H{"error": mfResp.Message})
			return
		}

		// --- Save to Database ---
		expiry, _ := time.Parse(time.RFC3339, mfResp.Data.SessionExpiry)

		session := models.Session{
			ID:                  sessionId,
			MyFatoorahSessionID: mfResp.Data.SessionId,
			EncryptionKey:       mfResp.Data.EncryptionKey,
			SessionExpiry:       expiry,
			Amount:              req.Amount,
			Currency:            "KWD",
			CustomerName:        customerName,
			CustomerEmail:       customerEmail,
			CustomerReference:   udf1Value, // UDF1 logic
		}

		if err := db.Create(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session to DB"})
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

func CreateSessionHandlerv2(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. SAFE EXTRACTION: Prevent nil pointer panic for OrderID
		orderRef := uuid.New().String()[:8]
		if req.OrderID != nil && *req.OrderID != "" {
			orderRef = *req.OrderID
		}

		// 2. CREATE INTERNAL INVOICE (The Parent record)
		invoice := models.Invoice{
			ID:            uuid.New(),
			OrderID:       orderRef,
			TotalValue:    req.Amount,
			Currency:      "KWD",
			Status:        "PENDING",
			CustomerName:  "Diaa Dawood",
			CustomerEmail: "diaadawood.mas@gmail.com",
			Source:        models.SourceEmbedded,
			Type:          models.TypeOneTime,
			CreatedAt:     time.Now(),
		}

		if err := db.Create(&invoice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create internal invoice"})
			return
		}

		// 3. PREPARE MYFATOORAH API CALL
		apiURL := "https://apitest.myfatoorah.com/v3/sessions"
		apiKey := "SK_KWT_Uw42KMHRKdVNZLHC4KEBTkHxOMfVsQi4txBkvQ1A1iAbVDdfPRhFsMQLVLkLXY3S"

		payload := map[string]interface{}{
			"PaymentMode": "COMPLETE_PAYMENT",
			"Order": map[string]interface{}{
				"Amount":             req.Amount,
				"Currency":           "KWD",
				"ExternalIdentifier": invoice.ID.String(), // Pass internal UUID
			},
			"MetaData": map[string]interface{}{
				"UDF1": invoice.ID.String(), // Anchor for Webhook to find this Invoice
			},
			"SupportedNetworks":       []string{"visa", "masterCard"},
			"SupportedPaymentMethods": []string{"card", "knet", "googlepay", "applepay"},
		}

		jsonData, _ := json.Marshal(payload)
		httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to gateway"})
			return
		}
		defer resp.Body.Close()

		// 4. PARSE GATEWAY RESPONSE
		var rawResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse gateway response"})
			return
		}

		if success, ok := rawResp["IsSuccess"].(bool); !ok || !success {
			msg := "Gateway error"
			if m, ok := rawResp["Message"].(string); ok {
				msg = m
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		data := rawResp["Data"].(map[string]interface{})

		// 5. SAVE SESSION RECORD (But NO Transaction yet)
		sessionUUID := uuid.New()

		mfExpiryStr, _ := data["SessionExpiry"].(string)
		mfExpiry, _ := time.Parse(time.RFC3339, mfExpiryStr)

		sessionRecord := models.Session{
			ID:                  sessionUUID,
			InvoiceID:           invoice.ID,
			MyFatoorahSessionID: data["SessionId"].(string),
			EncryptionKey:       data["EncryptionKey"].(string),
			SessionExpiry:       mfExpiry,
			Amount:              req.Amount,
			CreatedAt:           time.Now(),
		}

		if err := db.Create(&sessionRecord).Error; err != nil {
			// If session save fails, we don't necessarily need to rollback invoice,
			// but we should inform the client.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session record"})
			return
		}

		// 6. RETURN SUCCESS
		c.JSON(http.StatusCreated, gin.H{
			"message":             "Invoice created and session generated",
			"internal_invoice_id": invoice.ID,
			"session_id":          sessionRecord.ID,
			"session": gin.H{
				"id":             sessionRecord.MyFatoorahSessionID,
				"encryption_key": sessionRecord.EncryptionKey,
				"expiry":         sessionRecord.SessionExpiry,
			},
		})
	}
}
