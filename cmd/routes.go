package main

import (
	"github.com/aniqaqill/runners-list/handlers"
	"github.com/gofiber/fiber/v2"
)

// setupRoutes sets up the routes
func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/events", handlers.ListEvents)
	app.Post("/create-events", handlers.CreateEvents)
}
