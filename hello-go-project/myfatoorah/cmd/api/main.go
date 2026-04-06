package main

import (
	"fmt"
	"hello-go-project/myfatoorah/internal/database"
	"hello-go-project/myfatoorah/internal/models"

	_ "hello-go-project/myfatoorah/cmd/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MyFatoorah API
// @version 1.0
// @description This is a sample server for MyFatoorah.
// @host 2c03-196-156-78-10.ngrok-free.app
// @BasePath /

func main() {
	db := database.InitDB()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/invoices", models.GetAllInvoices(db))

	r.POST("/index.php", func(c *gin.Context) {
		// 2. Capture the "route" parameter from the URL
		route := c.Query("route")

		// 3. Check if it matches the MyFatoorah webhook path
		if route == "extension/myfatoorah/payment/myfatoorah.webhook" {

			// Handle the MyFatoorah Logic Here
			var webhookData map[string]interface{}
			if err := c.ShouldBindJSON(&webhookData); err != nil {
				c.JSON(400, gin.H{"error": "Invalid JSON"})
				return
			}

			fmt.Println("Webhook Received:", webhookData)
			c.JSON(200, gin.H{"status": "ok"})

		} else {
			// If it's index.php but the wrong route parameter
			c.JSON(404, gin.H{"message": "Route not found"})
		}
	})

	r.Run(":8080")
}
