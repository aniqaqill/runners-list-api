package main

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestMain(t *testing.T) {
	// Initialize the Fiber app
	app := fiber.New()

	// Define a simple route for testing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(":8080"); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Shutdown the server after a short delay
	defer func() {
		if err := app.Shutdown(); err != nil {
			t.Logf("Shutdown error: %v", err)
		}
	}()
}
