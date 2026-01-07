package http

import (
	"errors"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/core/service"
	"github.com/gofiber/fiber/v2"
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
	var user domain.Users

	// Parse the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	// Check for empty username or password
	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username and password cannot be empty",
		})
	}

	// Call the UserService to register the user
	if err := h.userService.Register(user.Username, user.Password); err != nil {
		if errors.Is(err, service.ErrUsernameAlreadyExists) {
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

// Handle login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var user domain.Users

	// Parse the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	// Call the UserService to authenticate the user
	authenticatedUser, err := h.userService.Login(user.Username, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to login",
		})
	}

	// Call the UserService to create the JWT token
	token, err := h.userService.CreateToken(int(authenticatedUser.ID))
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
		"token":   token,
	})
}

// ListUsers handles listing all users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve users",
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "No users found in the current database",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  users,
	})
}
