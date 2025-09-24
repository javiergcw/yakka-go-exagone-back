package models

import (
	"time"

	"github.com/google/uuid"
)

// JobRequirement represents job requirements in the system
type JobRequirement struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:100;not null;unique"`
	Description *string   `json:"description" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the JobRequirement model
func (JobRequirement) TableName() string {
	return "job_requirements"
}
