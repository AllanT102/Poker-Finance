package config

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"backend/internal/models"
    "os"
)

var DB *gorm.DB

func InitDatabase() {
    DB, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    // Enable UUID extension if not already enabled
    DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

    // AutoMigrate the schema
    err = DB.AutoMigrate(
        &models.User{},
        &models.Game{},
        &models.Transaction{},
        &models.PaymentDetails{},
        &models.PlayedGames{},
    )
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    log.Println("Database migration completed")
}