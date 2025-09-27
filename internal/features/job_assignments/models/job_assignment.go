package models

import (
	"time"

	"github.com/google/uuid"
)

// JobAssignment represents a job assignment in the system
type JobAssignment struct {
	ID            uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobID         uuid.UUID        `json:"job_id" gorm:"type:uuid;not null"`
	LabourUserID  uuid.UUID        `json:"labour_user_id" gorm:"type:uuid;not null"`
	ApplicationID uuid.UUID        `json:"application_id" gorm:"type:uuid;not null;uniqueIndex"`
	StartDate     *time.Time       `json:"start_date" gorm:"type:date"`
	EndDate       *time.Time       `json:"end_date" gorm:"type:date"`
	Status        AssignmentStatus `json:"status" gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	CreatedAt     time.Time        `json:"created_at" gorm:"not null;type:timestamptz;default:now()"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"not null;type:timestamptz;default:now()"`
}

// TableName returns the table name for the JobAssignment model
func (JobAssignment) TableName() string {
	return "job_assignments"
}
