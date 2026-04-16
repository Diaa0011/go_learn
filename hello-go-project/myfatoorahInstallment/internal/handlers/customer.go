package handlers

import (
	"hello-go-project/myfatoorahInstallment/internal/dto/request"
	"hello-go-project/myfatoorahInstallment/internal/dto/response"
	"hello-go-project/myfatoorahInstallment/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateCustomer godoc
// @Summary      Create or Update a Customer
// @Description  Accepts customer details and performs an upsert. Returns the customer ID and info.
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        customer  body      request.CreateCustomerRequest  true  "Customer Details"
// @Success      200       {object}  response.StandardResponse{data=response.CustomerResponse} "Success"
// @Failure      400       {object}  response.StandardResponse "Invalid Input"
// @Failure      500       {object}  response.StandardResponse "Database Error"
// @Router       /customers [post]
func CreateCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateCustomerRequest // Using your request struct

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, response.StandardResponse{
				Message: "Validation failed: " + err.Error(),
			})
			return
		}

		var customer models.Customer
		result := db.Where("email = ?", req.Email).First(&customer)

		if result.Error != nil {
			// Logic for New Customer
			customer = models.Customer{
				ID:    uuid.New(),
				Name:  req.Name,
				Email: req.Email,
				Phone: req.Phone,
			}
			if err := db.Create(&customer).Error; err != nil {
				c.JSON(500, response.StandardResponse{Message: "Database error"})
				return
			}
		} else {
			// Logic for Existing Customer (Update info)
			customer.Name = req.Name
			customer.Phone = req.Phone
			db.Save(&customer)
		}

		// Map model to the clean Response struct
		resp := response.CustomerResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			Phone:     customer.Phone,
			CreatedAt: customer.CreatedAt,
		}

		c.JSON(200, response.StandardResponse{
			Message: "Customer created successfully",
			Data:    resp,
		})
	}
}

// GetAllCustomers godoc
// @Summary      Get All Customers
// @Description  Retrieves customers with a global search (name/email/mobile) and dynamic sorting.
// @Tags         Customers
// @Param        page      query     int     false  "Page number"
// @Param        size     query     int     false  "Page size"
// @Param        search    query     string  false  "Search by name, email, or mobile"
// @Param        sort      query     string  false  "Sort column (default: created_at)"
// @Param        sort_dir  query     string  false  "Sort direction (asc/desc)"
// @Router       /customers/getall [get]
func GetAllCustomers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CustomerQueryRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.StandardResponse{Message: "Invalid query parameters"})
			return
		}

		var customers []models.Customer
		var totalRecords int64
		query := db.Model(&models.Customer{})

		// 1. Search Logic
		if req.Search != "" {
			s := "%" + req.Search + "%"
			query = query.Where("name ILIKE ? OR email ILIKE ? OR phone LIKE ?", s, s, s)
		}

		// 2. Count and Sort
		query.Count(&totalRecords)

		direction := "DESC"
		if req.SortDir == "asc" || req.SortDir == "ASC" {
			direction = "ASC"
		}

		// 3. Execute Pagination
		offset := (req.Page - 1) * req.Size
		if err := query.Offset(offset).Limit(req.Size).Order(req.Sort + " " + direction).Find(&customers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, response.StandardResponse{Message: "Database error"})
			return
		}

		// 4. Map to CustomerResponse DTO
		var customerData []response.CustomerResponse
		for _, cust := range customers {
			customerData = append(customerData, response.CustomerResponse{
				ID:        cust.ID,
				Name:      cust.Name,
				Email:     cust.Email,
				Phone:     cust.Phone,
				CreatedAt: cust.CreatedAt,
			})
		}

		// 5. Build the Paginated Response using your specific struct
		paginatedResp := response.PaginatedResponse[response.CustomerResponse]{
			Data: customerData,
			Meta: response.PaginationMetadata{
				TotalRecords: totalRecords,
				TotalPages:   int((totalRecords + int64(req.Size) - 1) / int64(req.Size)),
				CurrentPage:  req.Page,
				Size:         req.Size, // This maps to "index" in your JSON
			},
		}

		c.JSON(http.StatusOK, response.StandardResponse{
			Message: "Customers retrieved successfully",
			Data:    paginatedResp,
		})
	}
}
