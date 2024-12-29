package middleware

import (
	"time"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateCreateEventInput validates the input for creating an event
func ValidateCreateEventInput(c *fiber.Ctx) error {
	var input domain.Events
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	// Validate the input using the validator library
	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Validation failed",
			"details": err.Error(),
		})
	}

	// Validate date range (e.g., not in the past)
	const layout = "2006-01-02"
	eventDate, _ := time.Parse(layout, input.Date)
	if eventDate.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Date must be in the future.",
		})
	}

	return c.Next()
}
