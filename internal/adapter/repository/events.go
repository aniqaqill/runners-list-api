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

// Upsert inserts a new event or updates existing one based on name + date
func (r *GormEventRepository) Upsert(event *domain.Events) error {
	// First, try to find existing event by name and date
	var existing domain.Events
	err := r.db.Where("name = ? AND date = ?", event.Name, event.Date).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Event doesn't exist, create new one
		return r.db.Create(event).Error
	} else if err != nil {
		// Some other error occurred
		return err
	}

	// Event exists, update it
	event.ID = existing.ID
	event.CreatedAt = existing.CreatedAt
	return r.db.Save(event).Error
}

func (r *GormEventRepository) EventNameExists(name string) bool {
	var event domain.Events
	err := r.db.Where("name = ?", name).First(&event).Error
	return err == nil
}

// BulkUpsert inserts or updates multiple events in a single transaction
func (r *GormEventRepository) BulkUpsert(events []domain.Events) (inserted int, updated int, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		for i := range events {
			event := &events[i]

			// Try to find existing event by name and date
			var existing domain.Events
			findErr := tx.Where("name = ? AND date = ?", event.Name, event.Date).First(&existing).Error

			if findErr == gorm.ErrRecordNotFound {
				// Event doesn't exist, create new one
				if createErr := tx.Create(event).Error; createErr != nil {
					return createErr
				}
				inserted++
			} else if findErr != nil {
				// Some other error occurred
				return findErr
			} else {
				// Event exists, update it
				event.ID = existing.ID
				event.CreatedAt = existing.CreatedAt
				if saveErr := tx.Save(event).Error; saveErr != nil {
					return saveErr
				}
				updated++
			}
		}
		return nil
	})

	return inserted, updated, err
}

/* The repository layer implements the interfaces (ports) defined in the port layer.
It interacts with external systems (e.g., databases) and provides data to the service layer. */
