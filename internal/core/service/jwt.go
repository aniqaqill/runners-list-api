package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (s *UserService) CreateToken(username string) (string, error) {
	// Create the JWT claims, including the username and expiration time
	claims := jwt.MapClaims{
		"name": username,
		"exp":  time.Now().AddDate(0, 1, 0).Unix(), // Adds 1 month to the current time
	}

	// Create the JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
