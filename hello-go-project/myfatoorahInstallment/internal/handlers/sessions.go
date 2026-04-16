package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hello-go-project/myfatoorahInstallment/internal/dto/request"
	"hello-go-project/myfatoorahInstallment/internal/dto/response"
	"hello-go-project/myfatoorahInstallment/internal/models"
	"hello-go-project/myfatoorahInstallment/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateSessionHandler initializes a payment journey by contacting MyFatoorah.
// It generates "Session A" (InitiateSID) which the Frontend uses to render the secure card form.
//
// @Summary      Initiate Payment Session
// @Description  Calls MyFatoorah InitiateSession, stores the intent in PaymentSessions table, and returns Session A.
// @Tags         Sessions
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateSessionRequest  true  "Session Details"
// @Success      201      {object}  models.PaymentSession
// @Failure      400      {object}  map[string]string "Invalid Request"
// @Failure      500      {object}  map[string]string "API or Database Error"
// @Router       /sessions/initiate [post]
func CreateSessionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiBaseURL := os.Getenv("MYFATOORAH_BASE_URL")
		apiV2 := os.Getenv("MYFATOORAH_API_V2")
		apiIntitateSession := os.Getenv("MYFATOORAH_INTIATE_SESSION")
		apiKey := os.Getenv("MYFATOORAH_TOKEN")
		apiURL := apiBaseURL + apiV2 + apiIntitateSession

		payload := map[string]interface{}{
			"PaymentMethodId":    req.PaymentMethodId,
			"CustomerIdentifier": req.CustomerIdentifier,
			"IsRecurring":        true,
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

		// bodyBytes, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read response body"})
		// 	return
		// }

		// fmt.Println("--- DEBUG: MYFATOORAH RESPONSE ---")
		// fmt.Println(string(bodyBytes))
		// fmt.Println("-----------------------------------")

		var mfResponse response.CreateSessionResponse
		if err := json.NewDecoder(resp.Body).Decode(&mfResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		id, err := uuid.Parse(mfResponse.Data.SessionId)
		session := models.PaymentSession{
			ID:         id,
			CustomerID: req.CustomerIdentifier,
			CreatedAt:  time.Now(),
		}

		if err := db.Create(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session to DB"})
			return
		}

		c.JSON(http.StatusCreated, session)

	}
}

// UpdateSessionToTokenized godoc
// @Summary      Tokenize and Update Session
// @Description  Receives Session B from the frontend after card tokenization. Updates the existing record with the new ExecutionSID and amount.
// @Tags         Sessions
// @Accept       json
// @Produce      json
// @Param        request  body      request.ExecuteRequest  true  "Execution Details (Session B)"
// @Success      200      {object}  map[string]string       "message: Session updated successfully, status: TOKENIZED"
// @Failure      400      {object}  map[string]string       "error: Invalid request payload"
// @Failure      404      {object}  map[string]string       "error: Session not found"
// @Failure      500      {object}  map[string]string       "error: Failed to update session"
// @Router       /sessions/execute [post]
func UpdateSessionToTokenized(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.ExecuteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.StandardResponse{
				Message: "Invalid request payload",
			})
			return
		}

		// 1. Fetch Session & Preload Customer
		var session models.PaymentSession
		if err := db.Preload("Customer").Where("id = ?", req.OriginalSessionID).First(&session).Error; err != nil {
			c.JSON(http.StatusNotFound, response.StandardResponse{
				Message: "Session not found",
			})
			return
		}

		// 2. Prepare MyFatoorah Data
		countryCode, localNumber := utils.SplitPhoneAndCode(session.Customer.Phone)

		// 3. CREATE INTERNAL INVOICE (DRAFT)
		// We create this before the API call to have our own reference
		invoiceNumber := fmt.Sprintf("INV-%d-%s", time.Now().Unix(), session.Customer.Name[:3])
		newInvoice := models.Invoice{
			ID:            uuid.New(),
			CustomerID:    session.Customer.ID,
			InvoiceNumber: invoiceNumber,
			Amount:        req.InvoiceValue,
			TotalAmount:   req.InvoiceValue,
			Status:        "PENDING", // Becomes PAID via Webhook
			DueDate:       time.Now(),
		}

		if err := db.Create(&newInvoice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, response.StandardResponse{
				Message: "Failed to initialize internal invoice",
			})
			return
		}

		// 4. Prepare MyFatoorah Payload
		mfPayload := map[string]interface{}{
			"SessionId":          req.SessionId,
			"InvoiceValue":       req.InvoiceValue,
			"CustomerName":       session.Customer.Name,
			"DisplayCurrencyIso": "KWD",
			"MobileCountryCode":  countryCode,
			"CustomerMobile":     localNumber,
			"CustomerEmail":      session.Customer.Email,
			"CustomerReference":  session.ID.String(), // Link MyFatoorah to our Invoice UUID
			"Language":           "en",
			"ExternalIdentifier": session.ID.String(),
			"RecurringModel": map[string]interface{}{
				"RecurringType": req.RecurringModel.RecurringType,
				"IntervalDays":  req.RecurringModel.IntervalDays,
				"Iteration":     req.RecurringModel.Iteration,
				"RetryCount":    req.RecurringModel.RetryCount,
			},
		}

		// 5. Send Request to MyFatoorah
		apiBaseURL := os.Getenv("MYFATOORAH_BASE_URL")
		apiKey := os.Getenv("MYFATOORAH_TOKEN")
		apiURL := apiBaseURL + os.Getenv("MYFATOORAH_API_V2") + os.Getenv("MYFATOORAH_EXECUTE_SESSION")

		jsonData, _ := json.Marshal(mfPayload)
		httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.StandardResponse{Message: "Connection error"})
			return
		}
		defer resp.Body.Close()

		// 6. Decode Response
		var mfResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&mfResponse)

		// 7. Handle Success/Failure
		isSuccess, _ := mfResponse["IsSuccess"].(bool)
		if !isSuccess {
			// Mark internal invoice as failed if API fails immediately
			db.Model(&newInvoice).Update("Status", "FAILED")
			c.JSON(http.StatusBadRequest, response.StandardResponse{
				Message: "MyFatoorah Execution Failed",
				Data:    mfResponse,
			})
			return
		}

		// 8. Capture MyFatoorah Invoice ID and Update Local Records
		data, _ := mfResponse["Data"].(map[string]interface{})
		mfInvoiceID, _ := data["InvoiceId"].(float64)

		// Update our Invoice with the real MyFatoorah ID
		db.Model(&newInvoice).Updates(models.Invoice{
			MFInvoiceID: int(mfInvoiceID),
		})

		// Update Session
		session.ExecutionSID = req.SessionId
		session.Amount = req.InvoiceValue
		session.Status = "AWAITING_WEBHOOK"
		session.InvoiceID = int(mfInvoiceID) // Storing it in session for webhook lookup
		db.Save(&session)

		c.JSON(http.StatusOK, response.StandardResponse{
			Message: "Payment initiated. Awaiting OTP.",
			Data:    mfResponse,
		})
	}
}
