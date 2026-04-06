package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	OrderID           string `gorm:"uniqueIndex;not null"`
	InvoiceID         int    `gorm:"uniqueIndex"`
	InvoiceStatus     string `gorm:"type:varchar(20);default:'Pending'"`
	InvoiceReference  string
	InvoiceValue      float64 `gorm:"type:decimal(10,2)"`
	Currency          string  `gorm:"type:varchar(3)"`
	TransactionStatus string
}

// GetAllInvoices godoc
// @Summary      Get All Invoices
// @Tags         invoices
// @Success      200 {array} models.Invoice
// @Router       /invoices [get]
func GetAllInvoices(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var invoices []Invoice

		// SELECT * FROM invoices;
		result := db.Find(&invoices)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch invoices"})
			return
		}

		// Return the list as JSON
		c.JSON(http.StatusOK, invoices)
	}
}
