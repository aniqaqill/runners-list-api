package main

import (
	"log"

	"github.com/aniqaqill/runners-list/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Connect to the database
	database.ConnectDb()

	// Set up routes
	setupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
