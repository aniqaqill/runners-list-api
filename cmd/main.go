package main

import (
	"log"

	"github.com/aniqaqill/runners-list/internal/adapter/database"
	"github.com/aniqaqill/runners-list/internal/adapter/http"
	"github.com/aniqaqill/runners-list/internal/adapter/repository"
	"github.com/aniqaqill/runners-list/internal/core/service"
	"github.com/gofiber/fiber/v2"
	// "github.com/joho/godotenv"
	// Import the routes package
)

func main() {
	// Load environment variables
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

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

	// Initialize and start Fiber app
	if err := runFiberServer(eventHandler, userHandler); err != nil {
		log.Fatal(err)
	}
}

func runFiberServer(eventHandler *http.EventHandler, userHandler *http.UserHandler) error {
	app := fiber.New()

	// Call setupRoutes from the cmd package
	setupRoutes(app, eventHandler, userHandler)

	return app.Listen(":8080")
}
