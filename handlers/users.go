package handlers

import (
	"os"
	"time"

	"github.com/aniqaqill/runners-list/database"
	"github.com/aniqaqill/runners-list/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Register handler to create a new user
func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Username: data["username"],
		Password: string(password),
	}

	database.DB.Db.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "User registered successfully",
	})
}

// Login handler to generate JWT token
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	var user models.User
	database.DB.Db.Where("username = ?", data["username"]).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name": user.Username,
		"exp":  time.Now().AddDate(0, 1, 0).Unix(), // Adds 1 month to the current time
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Could not login",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Success",
		"token":   t,
	})
}
