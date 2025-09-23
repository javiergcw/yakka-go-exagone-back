package models

import (
	"time"

	"github.com/google/uuid"
)

// Jobsite represents a jobsite in the system
type Jobsite struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BuilderID   uuid.UUID `json:"builder_id" gorm:"type:uuid;not null;index"`
	Address     string    `json:"address" gorm:"type:text;not null"`
	City        string    `json:"city" gorm:"type:varchar(120);not null;index"`
	Suburb      *string   `json:"suburb" gorm:"type:varchar(120)"`
	Description *string   `json:"description" gorm:"type:text"`
	Latitude    float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude   float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Phone       *string   `json:"phone" gorm:"type:varchar(32)"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz;default:now()"`
}

// TableName returns the table name for the Jobsite model
func (Jobsite) TableName() string {
	return "jobsites"
}
