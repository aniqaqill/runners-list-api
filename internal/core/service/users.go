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
	repo      port.UserRepository
	jwtSecret string
}

// NewUserService creates a UserService. jwtSecret is injected here so that
// CreateToken does not read os.Getenv on each call.
func NewUserService(repo port.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(username, password string) error {
	existingUser, err := s.repo.FindByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		return ErrUsernameAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) ListUsers() ([]domain.Users, error) {
	return s.repo.FindAll()
}
