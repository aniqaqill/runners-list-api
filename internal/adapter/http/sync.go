package http

import (
	"time"

	"github.com/aniqaqill/runners-list/internal/core/domain"
)

// SyncRequest represents the request body for the internal sync endpoint
type SyncRequest struct {
	Events []SyncEventInput `json:"events" validate:"required,dive"`
}

// SyncEventInput represents a single event in the sync request
type SyncEventInput struct {
	Name            string `json:"name" validate:"required"`
	Location        string `json:"location" validate:"required"`
	State           string `json:"state"`
	Distance        string `json:"distance"`
	Date            string `json:"date" validate:"required"`
	Description     string `json:"description"`
	RegistrationURL string `json:"registration_url" validate:"required,url"`
}

// SyncResponse represents the response from the sync endpoint
type SyncResponse struct {
	Success  bool   `json:"success"`
	Inserted int    `json:"inserted"`
	Updated  int    `json:"updated"`
	Total    int    `json:"total"`
	Error    string `json:"error,omitempty"`
}

// ToEvent converts SyncEventInput to domain.Events
// Returns error if date parsing fails
func (s *SyncEventInput) ToEvent() (domain.Events, error) {
	// Parse date string to time.Time
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, s.Date)
	if err != nil {
		return domain.Events{}, err
	}

	return domain.Events{
		Name:             s.Name,
		Location:         s.Location,
		State:            s.State,
		Distance:         s.Distance,
		Date:             parsedDate,
		Description:      s.Description,
		RegisterationURL: s.RegistrationURL,
	}, nil
}
