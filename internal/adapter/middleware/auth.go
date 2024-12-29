package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

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
