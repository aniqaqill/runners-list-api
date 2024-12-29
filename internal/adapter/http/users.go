package http

import (
	"os"
	"time"

	"github.com/aniqaqill/runners-list/internal/core/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler with the given UserService
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var data map[string]string

	// Parse the request body into the data map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	// Call the UserService to register the user
	if err := h.userService.Register(data["username"], data["password"]); err != nil {
		if err.Error() == "username already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "Username already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to register user",
		})
	}

	// Return a success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "User registered successfully",
	})
}

// Login handles user login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parse the request body into the data map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	// Call the UserService to authenticate the user
	user, err := h.userService.Login(data["username"], data["password"])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	// Create the JWT claims, including the username and expiration time
	claims := jwt.MapClaims{
		"name": user.Username,
		"exp":  time.Now().AddDate(0, 1, 0).Unix(), // Adds 1 month to the current time
	}

	// Create the JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Could not login",
		})
	}

	// Return the JWT token in the response
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Success",
		"token":   t,
	})
}
