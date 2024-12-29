package port

import "github.com/aniqaqill/runners-list/internal/core/domain"

type EventRepository interface {
	Create(event *domain.Events) error
	FindAll() ([]domain.Events, error)
	FindByID(id uint) (*domain.Events, error)
	Delete(event *domain.Events) error
}
