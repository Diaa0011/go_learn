package database

import (
	"fmt"
	"hello-go-project/myfatoorah/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=Test2222!!!! dbname=myfatoorah_invoices port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("---> Running migrations...")
	err = db.AutoMigrate(&models.Invoice{})
	if err != nil {
		fmt.Println("!!! Migration Error:", err) // This will tell us EXACTLY why the table isn't there
	} else {
		fmt.Println("---> Migration successful!")
	}

	return db

}
