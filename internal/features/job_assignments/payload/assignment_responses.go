package payload

import (
	"time"

	"github.com/yakka-backend/internal/features/job_assignments/models"
)

// JobAssignmentResponse represents the response for a job assignment
type JobAssignmentResponse struct {
	ID            string                  `json:"id"`
	JobID         string                  `json:"job_id"`
	LabourUserID  string                  `json:"labour_user_id"`
	ApplicationID string                  `json:"application_id"`
	StartDate     *time.Time              `json:"start_date"`
	EndDate       *time.Time              `json:"end_date"`
	Status        models.AssignmentStatus `json:"status"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// CreateJobAssignmentResponse represents the response when creating a job assignment
type CreateJobAssignmentResponse struct {
	Assignment JobAssignmentResponse `json:"assignment"`
	Message    string                `json:"message"`
}

// UpdateJobAssignmentResponse represents the response when updating a job assignment
type UpdateJobAssignmentResponse struct {
	Assignment JobAssignmentResponse `json:"assignment"`
	Message    string                `json:"message"`
}

// GetJobAssignmentsResponse represents the response when getting job assignments
type GetJobAssignmentsResponse struct {
	Assignments []JobAssignmentResponse `json:"assignments"`
	Total       int64                   `json:"total"`
	Page        int                     `json:"page"`
	Limit       int                     `json:"limit"`
	TotalPages  int                     `json:"total_pages"`
}

// CompleteAssignmentResponse represents the response when completing an assignment
type CompleteAssignmentResponse struct {
	Assignment JobAssignmentResponse `json:"assignment"`
	Message    string                `json:"message"`
}

// CancelAssignmentResponse represents the response when cancelling an assignment
type CancelAssignmentResponse struct {
	Assignment JobAssignmentResponse `json:"assignment"`
	Message    string                `json:"message"`
}
