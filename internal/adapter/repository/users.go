package repository

import (
	"errors"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port"
	"gorm.io/gorm"
)

// GormUserRepository implements the UserRepository interface
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new instance of GormUserRepository
func NewGormUserRepository(db *gorm.DB) port.UserRepository {
	return &GormUserRepository{db: db}
}

// Save inserts a new user into the database
func (r *GormUserRepository) Create(user *domain.Users) error {
	return r.db.Create(user).Error
}

// FindByUsername retrieves a user by their username
func (r *GormUserRepository) FindByUsername(username string) (*domain.Users, error) {
	var user domain.Users
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if the user is not found
			return nil, nil
		}
		// Return the error for other cases
		return nil, err
	}
	return &user, nil
}

// Find All existing user
func (r *GormUserRepository) FindAll() ([]domain.Users, error) {
	var users []domain.Users
	err := r.db.Find(&users).Error
	return users, err
}
