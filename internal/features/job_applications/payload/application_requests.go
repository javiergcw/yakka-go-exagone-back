package payload

import "github.com/yakka-backend/internal/features/job_applications/models"

// CreateJobApplicationRequest represents the request to create a job application
type CreateJobApplicationRequest struct {
	JobID        string   `json:"job_id" validate:"required,uuid"`
	CoverLetter  *string  `json:"cover_letter"`
	ExpectedRate *float64 `json:"expected_rate" validate:"omitempty,min=0"`
	ResumeURL    *string  `json:"resume_url" validate:"omitempty,url"`
}

// UpdateJobApplicationRequest represents the request to update a job application
type UpdateJobApplicationRequest struct {
	Status       *models.ApplicationStatus `json:"status" validate:"omitempty,oneof=APPLIED REVIEWED ACCEPTED REJECTED WITHDRAWN"`
	CoverLetter  *string                   `json:"cover_letter"`
	ExpectedRate *float64                  `json:"expected_rate" validate:"omitempty,min=0"`
	ResumeURL    *string                   `json:"resume_url" validate:"omitempty,url"`
}

// GetJobApplicationsRequest represents the request to get job applications with filters
type GetJobApplicationsRequest struct {
	JobID        *string `json:"job_id" form:"job_id"`
	LabourUserID *string `json:"labour_user_id" form:"labour_user_id"`
	Status       *string `json:"status" form:"status"`
	Page         int     `json:"page" form:"page" validate:"min=1"`
	Limit        int     `json:"limit" form:"limit" validate:"min=1,max=100"`
}

// WithdrawApplicationRequest represents the request to withdraw an application
type WithdrawApplicationRequest struct {
	Reason *string `json:"reason"`
}
