package payload

// LabourApplicationRequest represents the request to apply for a job
type LabourApplicationRequest struct {
	JobID       string  `json:"job_id" validate:"required,uuid"`
	CoverLetter *string `json:"cover_letter" validate:"omitempty"`
	ResumeURL   *string `json:"resume_url" validate:"omitempty,url"`
}
