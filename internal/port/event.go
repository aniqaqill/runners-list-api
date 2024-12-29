package port

import "github.com/aniqaqill/runners-list/internal/core/domain"

type EventRepository interface {
	Save(event *domain.RunningEvents) error
	FindAll() ([]domain.RunningEvents, error)
	FindByID(id uint) (*domain.RunningEvents, error)
	Delete(event *domain.RunningEvents) error
}
