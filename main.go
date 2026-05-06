package main

import (
	"log"
	"os"
	"pubcore/database"
	"pubcore/router"

	"github.com/joho/godotenv"
)

// @title PubCore API
// @version 1.0
// @description Publishing platform API for managing Books, Authors, and Publishers.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@pubcore.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	// Get the PORT environment variable, with a fallback to a default value
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080" // Default port if not specified in .env
	}

	database.StartDB()
	router.StartServer().Run(":" + PORT)
}