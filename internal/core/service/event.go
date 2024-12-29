package service

import (
	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port"
)

type EventService struct {
	repo port.EventRepository
}

func NewEventService(repo port.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event *domain.RunningEvents) error {
	return s.repo.Save(event)
}

func (s *EventService) ListEvents() ([]domain.RunningEvents, error) {
	return s.repo.FindAll()
}

func (s *EventService) DeleteEvent(id uint) error {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(event)
}
