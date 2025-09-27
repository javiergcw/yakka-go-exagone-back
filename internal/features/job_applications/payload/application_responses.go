package payload

import (
	"time"

	"github.com/yakka-backend/internal/features/job_applications/models"
)

// JobApplicationResponse represents the response for a job application
type JobApplicationResponse struct {
	ID           string                   `json:"id"`
	JobID        string                   `json:"job_id"`
	LabourUserID string                   `json:"labour_user_id"`
	Status       models.ApplicationStatus `json:"status"`
	CoverLetter  *string                  `json:"cover_letter"`
	ExpectedRate *float64                 `json:"expected_rate"`
	ResumeURL    *string                  `json:"resume_url"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	WithdrawnAt  *time.Time               `json:"withdrawn_at"`
}

// CreateJobApplicationResponse represents the response when creating a job application
type CreateJobApplicationResponse struct {
	Application JobApplicationResponse `json:"application"`
	Message     string                 `json:"message"`
}

// UpdateJobApplicationResponse represents the response when updating a job application
type UpdateJobApplicationResponse struct {
	Application JobApplicationResponse `json:"application"`
	Message     string                 `json:"message"`
}

// GetJobApplicationsResponse represents the response when getting job applications
type GetJobApplicationsResponse struct {
	Applications []JobApplicationResponse `json:"applications"`
	Total        int64                    `json:"total"`
	Page         int                      `json:"page"`
	Limit        int                      `json:"limit"`
	TotalPages   int                      `json:"total_pages"`
}

// WithdrawApplicationResponse represents the response when withdrawing an application
type WithdrawApplicationResponse struct {
	Message string `json:"message"`
}
