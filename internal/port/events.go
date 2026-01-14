package port

import "github.com/aniqaqill/runners-list/internal/core/domain"

type EventRepository interface {
	Create(event *domain.Events) error
	FindAll() ([]domain.Events, error)
	FindByID(id uint) (*domain.Events, error)
	Delete(event *domain.Events) error
	EventNameExists(name string) bool
	Upsert(event *domain.Events) error
	BulkUpsert(events []domain.Events) (inserted int, updated int, err error)
}

/* The port layer defines the interfaces (ports) that the core layer
uses to interact with external systems (e.g., databases, HTTP clients). These interfaces are implemented by adapters. */
