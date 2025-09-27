package payload

import "time"

// LabourApplicantInfo represents the labour user information for an applicant
type LabourApplicantInfo struct {
	UserID    string  `json:"user_id"`
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url"`
	Phone     *string `json:"phone"`
	Email     string  `json:"email"`
}

// JobApplicantInfo represents a job application with labour information
type JobApplicantInfo struct {
	ApplicationID string              `json:"application_id"`
	Status        string              `json:"status"`
	CoverLetter   *string             `json:"cover_letter"`
	ExpectedRate  *float64            `json:"expected_rate"`
	ResumeURL     *string             `json:"resume_url"`
	AppliedAt     time.Time           `json:"applied_at"`
	Labour        LabourApplicantInfo `json:"labour"`
}

// JobWithApplicants represents a job with all its applicants
type JobWithApplicants struct {
	JobID      string             `json:"job_id"`
	JobTitle   string             `json:"job_title"`
	JobStatus  string             `json:"job_status"`
	CreatedAt  time.Time          `json:"created_at"`
	Applicants []JobApplicantInfo `json:"applicants"`
}

// BuilderApplicantsResponse represents the response for builder applicants
type BuilderApplicantsResponse struct {
	Jobs    []JobWithApplicants `json:"jobs"`
	Total   int                 `json:"total"`
	Message string              `json:"message"`
}

// BuilderApplicantDecisionResponse represents the response when hiring or rejecting an applicant
type BuilderApplicantDecisionResponse struct {
	ApplicationID string  `json:"application_id"`
	Hired         bool    `json:"hired"`
	AssignmentID  *string `json:"assignment_id,omitempty"` // Only present if hired
	Message       string  `json:"message"`
}
