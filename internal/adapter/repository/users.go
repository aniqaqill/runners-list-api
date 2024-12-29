package repository

// import (
// 	"github.com/aniqaqill/runners-list/internal/core/domain"
// 	"github.com/aniqaqill/runners-list/internal/port"
// 	"gorm.io/gorm"
// )

// // GormUserRepository implements the UserRepository interface
// type GormUserRepository struct {
// 	db *gorm.DB
// }

// // Delete removes a user from the database by their ID
// func (r *GormUserRepository) Delete(userID uint) error {
// 	return r.db.Delete(&domain.User{}, userID).Error
// }

// // NewGormUserRepository creates a new instance of GormUserRepository
// func NewGormUserRepository(db *gorm.DB) port.UserRepository {
// 	return &GormUserRepository{db: db}
// }

// // Save inserts a new user into the database
// func (r *GormUserRepository) Save(user *domain.User) error {
// 	return r.db.Create(user).Error
// }

// // FindByUsername retrieves a user by their username
// func (r *GormUserRepository) FindByUsername(username string) (*domain.User, error) {
// 	var user domain.User
// 	err := r.db.Where("username = ?", username).First(&user).Error
// 	return &user, err
// }
