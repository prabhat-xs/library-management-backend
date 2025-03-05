package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := os.Getenv("DATABASE_URL")

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	fmt.Println("Successfully connected to Neon Postgres database!")

}
