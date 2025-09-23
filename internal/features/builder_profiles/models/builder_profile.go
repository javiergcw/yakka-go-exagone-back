package models

import (
	"time"

	"github.com/google/uuid"
)

// BuilderProfile represents a builder profile in the system
type BuilderProfile struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	CompanyName *string   `json:"company_name" gorm:"size:255"`
	DisplayName *string   `json:"display_name" gorm:"size:255"`
	Location    *string   `json:"location" gorm:"size:255"`
	Bio         *string   `json:"bio" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the BuilderProfile model
func (BuilderProfile) TableName() string {
	return "builder_profiles"
}
