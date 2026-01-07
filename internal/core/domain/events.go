package domain

import (
	"time"

	"gorm.io/gorm"
)

// Events represents the running events data
type Events struct {
	gorm.Model
	Name             string    `json:"name" gorm:"type:text;not null" validate:"required"`
	Location         string    `json:"location" gorm:"type:text;not null" validate:"required"`
	State            string    `json:"state" gorm:"type:text"`
	Distance         string    `json:"distance" gorm:"type:text"`
	Date             time.Time `json:"date" gorm:"type:date;not null" validate:"required"`
	Description      string    `json:"description" gorm:"type:text" validate:"omitempty"`
	RegisterationURL string    `json:"registration_url" gorm:"type:text;not null" validate:"required,url"`
}

// The domain layer contains the core business logic and models. It defines the entities and interfaces (ports) that represent the business rules.
