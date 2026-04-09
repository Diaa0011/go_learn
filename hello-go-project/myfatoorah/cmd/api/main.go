package main

import (
	"hello-go-project/myfatoorah/internal/database"
	"hello-go-project/myfatoorah/internal/handlers"

	_ "hello-go-project/myfatoorah/cmd/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @title MyFatoorah API
// @version 1.0
// @description This is a sample server for MyFatoorah.
// @BasePath /

func main() {
	db := database.InitDB()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// New Session endpoint
	r.GET("/sessions", handlers.GetAllSessions(db))

	r.GET("/transactions", handlers.GetAllTransactions(db))
	r.POST("/sessions", handlers.CreateSessionHandler(db))

	r.POST("/index.php", MyFatoorahRouterHandler(db))

	r.Run(":8080")
}

// MyFatoorahRouterHandler handles the legacy PHP-style routing
// @Summary MyFatoorah Webhook
// @Description Receives status updates from MyFatoorah via the route query parameter.
// @Tags webhooks
// @Accept json
// @Produce json
// @Param route query string true "Must be: extension/myfatoorah/payment/myfatoorah.webhook"
// @Success 200 {string} string "OK"
// @Failure 404 {object} map[string]string "message: Route not found"
// @Router /index.php [post]
func MyFatoorahRouterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.Query("route")
		if route == "extension/myfatoorah/payment/myfatoorah.webhook" {
			handlers.MyFatoorahWebhookHandler(db)(c)
		} else {
			c.JSON(404, gin.H{"message": "Route not found"})
		}
	}
}
