package main

import (
	"github.com/aniqaqill/runners-list/handlers"
	"github.com/aniqaqill/runners-list/middleware"
	"github.com/gofiber/fiber/v2"
)

// setupRoutes sets up the routes
func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/events", handlers.ListEvents)
	app.Post("/create-events", middleware.ValidateCreateEventInput, handlers.CreateEvents)
	app.Delete("/events/:id", handlers.DeleteEvents)
}
