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
			c.JSON(http.StatusBadRequest, response.StandardResponse{Message: "Invalid payload"})
			return
		}

		var session models.PaymentSession
		if err := db.Preload("Customer").Where("id = ?", req.OriginalSessionID).First(&session).Error; err != nil {
			c.JSON(http.StatusNotFound, response.StandardResponse{Message: "Session not found"})
			return
		}
		installmentId := uuid.New()
		// 1. Prepare MyFatoorah Payload (Same as yours)
		countryCode, localNumber := utils.SplitPhoneAndCode(session.Customer.Phone)
		mfPayload := map[string]interface{}{
			"SessionId":          req.SessionId,
			"InvoiceValue":       req.InvoiceValue,
			"ExternalIdentifier": installmentId,
			"CustomerName":       session.Customer.Name,
			"DisplayCurrencyIso": "KWD",
			"MobileCountryCode":  countryCode,
			"CustomerMobile":     localNumber,
			"CustomerEmail":      session.Customer.Email,
			"CustomerReference":  session.Customer.ID,
			"Language":           "en",
			"RecurringModel": map[string]interface{}{
				"RecurringType": req.RecurringModel.RecurringType,
				"IntervalDays":  req.RecurringModel.IntervalDays,
				"Iteration":     req.RecurringModel.Iteration,
				"RetryCount":    req.RecurringModel.RetryCount,
			},
		}

		// 2. Call MyFatoorah
		apiURL := os.Getenv("MYFATOORAH_BASE_URL") + "/v2/ExecutePayment" // Simplified for example
		jsonData, _ := json.Marshal(mfPayload)
		httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("MYFATOORAH_TOKEN"))
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.StandardResponse{Message: "API Connection Error"})
			return
		}
		defer resp.Body.Close()

		var mfResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&mfResponse)

		if isSuccess, _ := mfResponse["IsSuccess"].(bool); !isSuccess {
			c.JSON(http.StatusBadRequest, response.StandardResponse{Message: "MF Execution Failed", Data: mfResponse})
			return
		}

		// 3. THE ATOMIC DRAFT (The core change)
		data := mfResponse["Data"].(map[string]interface{})
		mfInvoiceID := int(data["InvoiceId"].(float64))
		mfRecurringID := data["RecurringId"].(string)

		err = db.Transaction(func(tx *gorm.DB) error {
			// 1. Create the Draft Installment FIRST
			// We use the installmentId we generated at the top of the function
			newInstallment := models.Installment{
				ID:               installmentId, // This matches the ID you sent to MyFatoorah
				CustomerID:       session.CustomerID,
				MFRecurringID:    mfRecurringID,
				Status:           "PENDING",
				TotalAmount:      req.InvoiceValue,
				IterationAmount:  req.InvoiceValue,
				TotalIterations:  req.RecurringModel.Iteration,
				CurrentIteration: 0,
			}
			if err := tx.Create(&newInstallment).Error; err != nil {
				return err
			}

			// 2. Create the Draft Invoice and LINK IT
			newInvoice := models.Invoice{
				ID:            uuid.New(),
				CustomerID:    session.CustomerID,
				InstallmentID: installmentId, // CRUCIAL: Link to the installment above
				MFInvoiceID:   mfInvoiceID,
				InvoiceNumber: fmt.Sprintf("INV-%d", time.Now().Unix()),
				Amount:        req.InvoiceValue,
				TotalAmount:   req.InvoiceValue,
				Status:        "PENDING",
				DueDate:       time.Now(),
			}
			if err := tx.Create(&newInvoice).Error; err != nil {
				return err
			}

			// 3. Update Session
			session.Status = "AWAITING_WEBHOOK"
			session.InvoiceID = mfInvoiceID
			session.ExecutionSID = req.SessionId
			// Note: If session belongs to the same DB transaction, use tx.Save
			if err := tx.Save(&session).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, response.StandardResponse{Message: "DB Error during drafting"})
			return
		}

		c.JSON(http.StatusOK, response.StandardResponse{Message: "Draft created. Awaiting payment.", Data: mfResponse})
	}
}
