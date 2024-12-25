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

	// User registration and login routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Apply JWT middleware to protect POST and DELETE routes
	app.Post("/create-events", middleware.JWTProtected(), middleware.ValidateCreateEventInput, handlers.CreateEvents)
	app.Delete("/events/:id", middleware.JWTProtected(), handlers.DeleteEvents)
}
