package models

import (
	"time"

	"github.com/google/uuid"
)

// LabourProfile represents a labour profile in the system
type LabourProfile struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Location  *string   `json:"location" gorm:"size:255"`
	Bio       *string   `json:"bio" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the LabourProfile model
func (LabourProfile) TableName() string {
	return "labour_profiles"
}
