package config

import (
	"fmt"
	"log"
	"os"

	"github.com/greatdaveo/privycode-server/internal/models"
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
			log.Fatal("❌ Error loading .env file")
		}
	}

	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error connecting to DB: %v", err)
	}

	// To verify the connection pinging
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("❌ Error getting raw DB handle: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Error pingin DB: %v", err)
	}

	DB = database

	fmt.Println("Connected to PostgreSQL successfully!!! ✅")

}

func RunMigrations() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate: %v", err)
	}

	fmt.Println("Migrations completed successfully ✅")

}
