package service

// import (
// 	"github.com/aniqaqill/runners-list/internal/core/domain"
// 	"github.com/aniqaqill/runners-list/internal/port"
// 	"golang.org/x/crypto/bcrypt"
// )

// type UserService struct {
// 	repo port.UserRepository
// }

// func NewUserService(repo port.UserRepository) *UserService {
// 	return &UserService{repo: repo}
// }

// func (s *UserService) Register(username, password string) error {
// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	user := &domain.User{
// 		Username: username,
// 		Password: string(hashedPassword),
// 	}
// 	return s.repo.Save(user)
// }

// func (s *UserService) Login(username, password string) (*domain.User, error) {
// 	user, err := s.repo.FindByUsername(username)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }
