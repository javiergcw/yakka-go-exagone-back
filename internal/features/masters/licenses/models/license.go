package models

import (
	"time"

	"github.com/google/uuid"
)

// License represents a license in the system
type License struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the License model
func (License) TableName() string {
	return "licenses"
}
