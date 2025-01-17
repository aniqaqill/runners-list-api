package port

import "github.com/aniqaqill/runners-list/internal/core/domain"

type UserRepository interface {
	Create(user *domain.Users) error
	FindByUsername(username string) (*domain.Users, error)
	FindAll() ([]domain.Users, error)
}
