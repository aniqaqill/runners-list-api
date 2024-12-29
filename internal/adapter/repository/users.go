package repository

import (
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
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
