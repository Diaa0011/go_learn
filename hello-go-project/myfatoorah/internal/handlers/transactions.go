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
