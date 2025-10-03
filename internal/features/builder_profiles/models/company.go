package models

import (
	"time"

	"github.com/google/uuid"
)

// Company represents a company in the system
type Company struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:255;not null;uniqueIndex"`
	Description *string   `json:"description" gorm:"type:text"`
	Website     *string   `json:"website" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the Company model
func (Company) TableName() string {
	return "companies"
}
