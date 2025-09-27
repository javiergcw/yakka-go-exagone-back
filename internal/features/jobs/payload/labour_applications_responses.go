package payload

import "time"

// LabourApplicationResponse represents the response when applying for a job
type LabourApplicationResponse struct {
	ApplicationID string    `json:"application_id"`
	JobID         string    `json:"job_id"`
	JobTitle      string    `json:"job_title"`
	Status        string    `json:"status"`
	CoverLetter   *string   `json:"cover_letter"`
	ResumeURL     *string   `json:"resume_url"`
	AppliedAt     time.Time `json:"applied_at"`
	Message       string    `json:"message"`
}
