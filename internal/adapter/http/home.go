package http

import (
	"github.com/gofiber/fiber/v2"
)

// func for index
func Home(c *fiber.Ctx) error {
	return c.SendString("Running App!")
}
