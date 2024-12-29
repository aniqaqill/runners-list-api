package domain

import "gorm.io/gorm"

// models for the running events data
type RunningEvents struct {
	gorm.Model
	Name             string `json:"name" gorm:"type:text;not null"`
	Location         string `json:"location" gorm:"type:text;not null"`
	Date             string `json:"date" gorm:"type:text;not null"`
	Description      string `json:"description" gorm:"type:text;not null"`
	RegisterationURL string `json:"registration_url" gorm:"type:text;not null"`
}
