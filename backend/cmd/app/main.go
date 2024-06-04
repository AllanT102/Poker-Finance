package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/services/email"
	"log"

	_ "backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Poker Finance API
// @version         1.0
// @description     An API used to manage payments with friends for home poker games :)

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	// Get database connection details from environment variables
	config.InitDatabase()

	// Initialize email queue
	email.CreateEmailChannel()

	router := gin.Default()
	api.SetupRoutes(router)

	log.Fatal(router.Run(":8080"))
}
