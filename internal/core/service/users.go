package service

import (
	"errors"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password string) error {
	// Check if the username already exists
	existingUser, _ := s.repo.FindByUsername(username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user := &domain.Users{
		Username: username,
		Password: string(hashedPassword),
	}

	// Create the user
	return s.repo.Create(user)
}

func (s *UserService) GetUserByUsername(username string) (*domain.Users, error) {
	return s.repo.FindByUsername(username)
}

func (s *UserService) Login(username, password string) (*domain.Users, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
