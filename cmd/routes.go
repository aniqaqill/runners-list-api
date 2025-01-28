package main

import (
	"github.com/aniqaqill/runners-list/internal/adapter/http"
	"github.com/aniqaqill/runners-list/internal/adapter/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App, eventHandler *http.EventHandler, userHandler *http.UserHandler) {
	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is working!")
	})

	app.Get("/events", eventHandler.ListEvents)
	app.Get("/users", userHandler.ListUsers)

	// User registration and login routes
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	// Protected routes (JWT middleware)
	protected := app.Group("/protected", middleware.JWTProtected())

	protected.Post("/events/create-events", middleware.ValidateCreateEventInput, eventHandler.CreateEvent)
	protected.Delete("/events/:id", eventHandler.DeleteEvent)
}
