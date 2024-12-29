package main

import (
	"log"

	"github.com/aniqaqill/runners-list/internal/adapter/database"
	"github.com/aniqaqill/runners-list/internal/adapter/http"
	"github.com/aniqaqill/runners-list/internal/adapter/repository"
	"github.com/aniqaqill/runners-list/internal/core/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	database.ConnectDb()

	// Initialize repositories
	eventRepo := repository.NewGormEventRepository(database.DB.Db)
	userRepo := repository.NewGormUserRepository(database.DB.Db)

	// Initialize services
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	eventHandler := http.NewEventHandler(eventService)
	userHandler := http.NewUserHandler(userService)

	// Initialize Fiber app
	app := fiber.New()

	// Set up routes
	setupRoutes(app, eventHandler, userHandler)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
