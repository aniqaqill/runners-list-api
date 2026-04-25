package port

import (
	"time"

	"github.com/aniqaqill/runners-list/internal/core/domain"
)

// EventFilter holds optional query parameters for listing events.
// Zero values mean "no filter" for that field.
type EventFilter struct {
	State  string    // filter by state name (case-insensitive)
	From   time.Time // events on or after this date
	To     time.Time // events on or before this date
	Limit  int       // page size (default 50)
	Offset int       // page offset (default 0)
}

type EventRepository interface {
	Create(event *domain.Events) error
	FindAll(filter EventFilter) ([]domain.Events, error)
	FindByID(id uint) (*domain.Events, error)
	Delete(event *domain.Events) error
	EventNameExists(name string) bool
	Upsert(event *domain.Events) error
	BulkUpsert(events []domain.Events) (inserted int, updated int, err error)
}
