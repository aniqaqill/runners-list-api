package middleware

import (
	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ValidateRegisterInput validates the input for user registration.
func ValidateRegisterInput(c *fiber.Ctx) error {
	var input domain.Users

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid input format",
		})
	}

	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "validation failed",
			"details": err.Error(),
		})
	}

	c.Locals("registerInput", input)
	return c.Next()
}

// ValidateLoginInput validates the input for user login.
func ValidateLoginInput(c *fiber.Ctx) error {
	var input domain.Users

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid input format",
		})
	}

	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "validation failed",
			"details": err.Error(),
		})
	}

	c.Locals("loginInput", input)
	return c.Next()
}

// JWTProtected returns a Fiber middleware that validates a Bearer JWT.
//
// The secret is injected at construction time (from Config) rather than read
// from os.Getenv on each request. This keeps middleware pure and testable.
func JWTProtected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if len(authHeader) < len("Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "unauthorized",
			})
		}

		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "unauthorized",
			})
		}

		return c.Next()
	}
}
