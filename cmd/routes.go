package main

import (
	"github.com/aniqaqill/runners-list/internal/adapter/http"
	"github.com/aniqaqill/runners-list/internal/adapter/middleware"
	"github.com/gofiber/fiber/v2"
)

// func setupRoutes(app *fiber.App, eventHandler *http.EventHandler, userHandler *http.UserHandler) {
func setupRoutes(app *fiber.App, eventHandler *http.EventHandler) {
	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Running App!")
	})
	app.Get("/events", eventHandler.ListEvents)

	// User registration and login routes
	// app.Post("/register", userHandler.Register)
	// app.Post("/login", userHandler.Login)

	// Protected routes (JWT middleware)
	app.Post("/create-events", middleware.JWTProtected(), middleware.ValidateCreateEventInput, eventHandler.CreateEvent)
	app.Delete("/events/:id", middleware.JWTProtected(), eventHandler.DeleteEvent)
}
