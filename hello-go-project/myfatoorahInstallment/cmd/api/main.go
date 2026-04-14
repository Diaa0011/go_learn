package main

import (
	"hello-go-project/myfatoorahInstallment/internal/database"
	"hello-go-project/myfatoorahInstallment/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MyFatoorah Installment API
// @version 1.0
// @description This is a sample server for MyFatoorah.
// @BasePath /

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file - make sure it exists in the root folder")
	}

	db := database.InitDB()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/sessions", handlers.CreateSessionHandler(db))
}
