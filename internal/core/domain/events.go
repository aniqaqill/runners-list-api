package domain

import (
	"time"

	"gorm.io/gorm"
)

// Events is the core domain entity for a running event.
//
// Note on gorm tag `column:registeration_url`:
//   The original Go field was misspelled as "RegisterationURL". Renaming it
//   here to the correct "RegistrationURL" would change the GORM-derived column
//   name from "registeration_url" → "registration_url", breaking existing data.
//   The explicit column tag pins the DB column to its current name so no
//   emergency data migration is needed when this code is deployed.
//   Run the migration in 001_rename_registration_url.sql when convenient to
//   clean up the DB column name as well.
type Events struct {
	gorm.Model
	Name            string    `json:"name"             gorm:"type:text;not null"                              validate:"required"`
	Location        string    `json:"location"         gorm:"type:text;not null"                              validate:"required"`
	State           string    `json:"state"            gorm:"type:text"`
	Distance        string    `json:"distance"         gorm:"type:text"`
	Date            time.Time `json:"date"             gorm:"type:date;not null;index"                        validate:"required"`
	Description     string    `json:"description"      gorm:"type:text"                                       validate:"omitempty"`
	RegistrationURL string    `json:"registration_url" gorm:"column:registeration_url;type:text;not null"    validate:"required,url"`
}
