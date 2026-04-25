package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CreateToken generates a signed JWT for the given user ID.
//
// The signing secret is now stored on UserService (injected at construction)
// rather than read from os.Getenv on every call. This makes the function
// deterministic and unit-testable without environment setup.
func (s *UserService) CreateToken(ID int) (string, error) {
	if ID <= 0 {
		return "", errors.New("ID cannot be negative or zero")
	}

	claims := jwt.MapClaims{
		"id":  ID,
		"exp": time.Now().AddDate(0, 1, 0).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("failed to sign the token")
	}

	return tokenString, nil
}
