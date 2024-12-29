package port

import "github.com/aniqaqill/runners-list/internal/core/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	FindUserByUsername(username string) (*domain.User, error)
}
