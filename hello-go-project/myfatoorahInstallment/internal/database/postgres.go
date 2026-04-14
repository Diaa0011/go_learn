package database

import (
	"fmt"
	"hello-go-project/myfatoorahInstallment/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	// dsn1 := "host=localhost user=postgres password=Test2222!!!! dbname=myfatoorah_installment port=5432 sslmode=disable"

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	dsnRoot := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode)

	rootDb, err := gorm.Open(postgres.Open(dsnRoot), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	var exists int
	rootDb.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbName).Scan(&exists)
	if exists == 0 {
		fmt.Printf("---> Database %s not found. Creating it now...\n", dbName)
		// We use Exec because CREATE DATABASE cannot run inside a transaction
		if err := rootDb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
			log.Fatal("Failed to create database: ", err)
		}
	}

	// Close the root connection
	sqlRoot, _ := rootDb.DB()
	sqlRoot.Close()

	// 2. Now connect to the actual project database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to project database: ", err)
	}

	// Connection Pool Settings
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("---> Running migrations...")
	err = db.AutoMigrate(&models.Installment{}, &models.Transaction{}, &models.PaymentSession{})
	if err != nil {
		fmt.Println("!!! Migration Error:", err)
	}

	return db
}
