package handlers

import "github.com/gofiber/fiber/v2"

// func for event creation
func CreateEvents(c *fiber.Ctx) error {
	return c.SendString("Event created")
}
