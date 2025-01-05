package domain

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" gorm:"not null" validate:"required"`
}
