package models

import (
	"time"

	"github.com/google/uuid"
)

// SportsQualification represents a sport in the qualifications system
type SportsQualification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"size:150;not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the SportsQualification model
func (SportsQualification) TableName() string {
	return "qualifications_sport"
}
