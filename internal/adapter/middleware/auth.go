package middleware

import (
	"os"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ValidateRegisterInput validates the input for user registration
func ValidateRegisterInput(c *fiber.Ctx) error {
	var input domain.Users

	// Parse the request body
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

	// Store the validated input in the context for use in the handler
	c.Locals("registerInput", input)

	return c.Next()
}

// ValidateLoginInput validates the input for user login
func ValidateLoginInput(c *fiber.Ctx) error {
	var input domain.Users

	// Parse the request body
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

	// Store the validated input in the context for use in the handler
	c.Locals("loginInput", input)

	return c.Next()
}

// JWTProtected protects routes with JWT authentication
func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the JWT from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Unauthorized",
			})
		}

		// Extract the token from the "Bearer " prefix
		tokenString := authHeader[len("Bearer "):]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Unauthorized",
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Unauthorized",
			})
		}

		return c.Next()
	}
}
