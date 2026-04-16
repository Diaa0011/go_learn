package database

import (
	"fmt"
	"hello-go-project/myfatoorahInstallment/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// dsn1 := "host=localhost user=postgres password=Test2222!!!! dbname=myfatoorah_installment port=5432 sslmode=disable"

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("---> Running migrations...")
	log.Println("Running migrations for Installments, Sessions, and Transactions...")

	err = db.AutoMigrate(&models.Customer{}, &models.Installment{}, &models.Transaction{}, &models.PaymentSession{}, &models.Invoice{})
	if err != nil {
		fmt.Println("!!! Migration Error:", err)
	}

	return db
}
