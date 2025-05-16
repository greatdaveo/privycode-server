package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	// To only load .env file in local development
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dsn := os.Getenv("DATABASE_URL")

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("❌ Error connecting to DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully!!! ✅")

}

func RunMigration() {}
