package handlers

import (
	"hello-go-project/myfatoorah/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MyFatoorahWebhookHandler godoc
// @Summary      MyFatoorah Webhook
// @Description  Receives status updates from MyFatoorah.
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Success      200 {string} string "OK"
// @Router       /index.php [post]
func MyFatoorahWebhookHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload MyFatoorahWebhook

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		result := db.Model(&models.Invoice{}).Where("invoice_id = ?", payload.Data.InvoiceId).Updates(models.Invoice{
			InvoiceStatus: payload.Data.InvoiceStatus,
		})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		c.Status(http.StatusOK)
	}
}

type MyFatoorahWebhook struct {
	Event string `json:"Event"` // e.g., "TransactionsStatusChanged"
	Data  struct {
		InvoiceId         int     `json:"InvoiceId"`
		InvoiceStatus     string  `json:"InvoiceStatus"` // e.g., "Paid"
		InvoiceValue      float64 `json:"InvoiceValue"`
		CustomerReference string  `json:"CustomerReference"` // Your internal OrderID
	} `json:"Data"`
}
