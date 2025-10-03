package models

import (
	"time"

	"github.com/google/uuid"
)

// Qualification represents a qualification/certification in the system
type Qualification struct {
	ID           uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SportID      uuid.UUID            `json:"sport_id" gorm:"type:uuid;not null;index"`
	Sport        *SportsQualification `json:"sport,omitempty" gorm:"foreignKey:SportID"`
	Title        string               `json:"title" gorm:"size:255;not null"`
	Organization *string              `json:"organization" gorm:"size:150"`
	Country      *string              `json:"country" gorm:"size:100"`
	Status       string               `json:"status" gorm:"size:50;default:active"`
	CreatedAt    time.Time            `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt    time.Time            `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the Qualification model
func (Qualification) TableName() string {
	return "qualifications"
}
