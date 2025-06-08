package main

import (
    "log"

    "github.com/aniqaqill/runners-list/internal/adapter/http"
    "github.com/aniqaqill/runners-list/internal/adapter/middleware"
    "github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App, eventHandler *http.EventHandler, userHandler *http.UserHandler) {
    // Group all routes under /api/v1
    api := app.Group("/api")
    v1 := api.Group("/v1")

    log.Println("setupRoutes is being called")
    v1.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("API v1 is working!")
    })

    v1.Get("/events", eventHandler.ListEvents)
    v1.Get("/users", userHandler.ListUsers)

    // User registration and login routes
    v1.Post("/register", userHandler.Register)
    v1.Post("/login", userHandler.Login)

    // Protected routes (JWT middleware)
    protected := v1.Group("/protected", middleware.JWTProtected())
    protected.Post("/events/create-events", middleware.ValidateCreateEventInput, eventHandler.CreateEvent)
    protected.Delete("/events/:id", eventHandler.DeleteEvent)
}