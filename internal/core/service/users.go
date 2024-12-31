package service

import (
	"errors"
	"fmt"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password string) error {
	// Check if the username already exists
	existingUser, err := s.repo.FindByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Return an error if the query fails (other than "record not found")
		return fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		// Return a conflict error if the username already exists
		return ErrUsernameAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the new user
	newUser := &domain.Users{
		Username: username,
		Password: string(hashedPassword),
	}
	if err := s.repo.Create(newUser); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *UserService) GetUserByUsername(username string) (*domain.Users, error) {
	return s.repo.FindByUsername(username)
}

func (s *UserService) Login(username, password string) (*domain.Users, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
