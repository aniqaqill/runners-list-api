package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// InternalAPIKeyAuth validates the internal API key for scraper access
func InternalAPIKeyAuth(c *fiber.Ctx) error {
	apiKey := c.Get("X-Internal-Token")
	expectedKey := os.Getenv("INTERNAL_API_KEY")

	if expectedKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal API key not configured",
		})
	}

	if apiKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Missing X-Internal-Token header",
		})
	}

	if apiKey != expectedKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid API key",
		})
	}

	return c.Next()
}
