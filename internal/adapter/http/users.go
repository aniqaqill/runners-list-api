package http

// import (
// 	"os"
// 	"time"

// 	"github.com/aniqaqill/runners-list/internal/core/service"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v4"
// )

// type UserHandler struct {
// 	userService *service.UserService
// }

// func NewUserHandler(userService *service.UserService) *UserHandler {
// 	return &UserHandler{userService: userService}
// }

// func (h *UserHandler) Register(c *fiber.Ctx) error {
// 	var data map[string]string

// 	if err := c.BodyParser(&data); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Invalid input format",
// 		})
// 	}

// 	if err := h.userService.Register(data["username"], data["password"]); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Failed to register user",
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"error":   false,
// 		"message": "User registered successfully",
// 	})
// }

// func (h *UserHandler) Login(c *fiber.Ctx) error {
// 	var data map[string]string

// 	if err := c.BodyParser(&data); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Invalid input format",
// 		})
// 	}

// 	user, err := h.userService.Login(data["username"], data["password"])
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Invalid credentials",
// 		})
// 	}

// 	// Create the Claims
// 	claims := jwt.MapClaims{
// 		"name": user.Username,
// 		"exp":  time.Now().AddDate(0, 1, 0).Unix(), // Adds 1 month to the current time
// 	}

// 	// Create token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Could not login",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"error":   false,
// 		"message": "Success",
// 		"token":   t,
// 	})
// }
