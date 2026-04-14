package handlers

import (
	"bytes"
	"encoding/json"
	"hello-go-project/myfatoorahInstallment/internal/dto/request"
	"hello-go-project/myfatoorahInstallment/internal/dto/response"
	"hello-go-project/myfatoorahInstallment/internal/models"
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
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateSessionRequest  true  "Session Details"
// @Success      201      {object}  models.PaymentSession
// @Failure      400      {object}  map[string]string "Invalid Request"
// @Failure      500      {object}  map[string]string "API or Database Error"
// @Router       /payments/initiate [post]
func CreateSessionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiURL := os.Getenv("MYFATOORAH_BASE_URL")
		apiKey := os.Getenv("MYFATOORAH_TOKEN")

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

		var mfResponse response.CreateSessionResponse
		if err := json.NewDecoder(resp.Body).Decode(&mfResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		session := models.PaymentSession{
			ID:          uuid.New(),
			UserID:      req.CustomerIdentifier,
			InitiateSID: mfResponse.Data.SessionId,
			CreatedAt:   time.Now(),
		}

		if err := db.Create(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session to DB"})
			return
		}

		c.JSON(http.StatusCreated, session)

	}
}
