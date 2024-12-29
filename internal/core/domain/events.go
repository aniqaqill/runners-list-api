package domain

import "gorm.io/gorm"

// Events represents the running events data
type Events struct {
	gorm.Model
	Name             string `json:"name" gorm:"type:text;not null" validate:"required"`
	Location         string `json:"location" gorm:"type:text;not null" validate:"required"`
	Date             string `json:"date" gorm:"type:text;not null" validate:"required,datetime=2006-01-02"`
	Description      string `json:"description" gorm:"type:text" validate:"omitempty"`
	RegisterationURL string `json:"registration_url" gorm:"type:text;not null" validate:"required,url"`
}

// The domain layer contains the core business logic and models. It defines the entities and interfaces (ports) that represent the business rules.
