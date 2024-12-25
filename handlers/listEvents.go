package handlers

import "github.com/gofiber/fiber/v2"

// func for index
func Home(c *fiber.Ctx) error {
	return c.SendString("Running App!")
}

// func for retrieving running event list
func ListEvents(c *fiber.Ctx) error {
	return c.SendString("Running Events List!")
}
