package models

import (
	"time"

	"github.com/google/uuid"
)

// JobApplication represents a job application in the system
type JobApplication struct {
	ID           uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobID        uuid.UUID         `json:"job_id" gorm:"type:uuid;not null"`
	LabourUserID uuid.UUID         `json:"labour_user_id" gorm:"type:uuid;not null"`
	Status       ApplicationStatus `json:"status" gorm:"type:varchar(20);not null;default:'APPLIED'"`
	CoverLetter  *string           `json:"cover_letter" gorm:"type:text"`
	ExpectedRate *float64          `json:"expected_rate" gorm:"type:decimal(12,2)"`
	ResumeURL    *string           `json:"resume_url" gorm:"type:text"`
	CreatedAt    time.Time         `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"not null;type:timestamptz"`
	WithdrawnAt  *time.Time        `json:"withdrawn_at" gorm:"type:timestamptz"`
}

// TableName returns the table name for the JobApplication model
func (JobApplication) TableName() string {
	return "job_applications"
}
