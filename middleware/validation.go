package middleware

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type EventInput struct {
	Name             string `json:"name" validate:"required"`
	Location         string `json:"location" validate:"required"`
	Date             string `json:"date" validate:"required,datetime=2006-01-02"`
	Description      string `json:"description" validate:"required"`
	RegisterationURL string `json:"registration_url" validate:"required,url"`
}

// ValidateCreateEventInput validates the input for creating an event
func ValidateCreateEventInput(c *fiber.Ctx) error {
	var input EventInput
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
