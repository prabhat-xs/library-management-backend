package config

import (
	"fmt"
	"log"
	"os"

	"github.com/prabhat-xs/library-management-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := os.Getenv("DATABASE_URL")

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		 Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	// Get the underlying SQL DB object
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to get DB object: %v", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Library{},
		&models.User{},
		&models.Books{},
		&models.RequestEvents{},
		&models.IssueRegistry{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Successfully connected to Neon Postgres database!")

}
