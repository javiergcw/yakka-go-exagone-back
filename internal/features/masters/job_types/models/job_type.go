package models

import (
	"time"

	"github.com/google/uuid"
)

// JobType represents job types in the system
type JobType struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:100;not null;unique"`
	Description *string   `json:"description" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the JobType model
func (JobType) TableName() string {
	return "job_types"
}
