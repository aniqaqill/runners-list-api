package repository

import (
	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port"
	"gorm.io/gorm"
)

type GormEventRepository struct {
	db *gorm.DB
}

func NewGormEventRepository(db *gorm.DB) port.EventRepository {
	return &GormEventRepository{db: db}
}

func (r *GormEventRepository) Create(event *domain.Events) error {
	return r.db.Create(event).Error
}

func (r *GormEventRepository) FindAll() ([]domain.Events, error) {
	var events []domain.Events
	err := r.db.Find(&events).Error
	return events, err
}

func (r *GormEventRepository) FindByID(id uint) (*domain.Events, error) {
	var event domain.Events
	err := r.db.First(&event, id).Error
	return &event, err
}

func (r *GormEventRepository) Delete(event *domain.Events) error {
	return r.db.Delete(event).Error
}