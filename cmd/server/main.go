package main

import (
	"log"
	"net/http"
	"os"

	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/middleware"
	"github.com/greatdaveo/privycode-server/internal/routes"
)

func main() {

	// To Connect to the DB
	config.ConnectDB()

	// For auto migrate to DB
	// config.RunMigrations()

	// To set up HTTP router
	mux := http.NewServeMux()

	routes.APIRoutes(mux)

	handlerWithCORS := middleware.WithCORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s... ✅", port)

	err := http.ListenAndServe(":"+port, handlerWithCORS)

	if err != nil {
		log.Fatalf("❌ Could not start sever: %v", err)
	}
}
