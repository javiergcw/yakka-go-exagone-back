package payload

import "time"

// LabourJobInfo represents a job with application status for a labour user
type LabourJobInfo struct {
	JobID           string     `json:"job_id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Location        string     `json:"location"`
	JobType         string     `json:"job_type"`
	ExperienceLevel string     `json:"experience_level"`
	Status          string     `json:"status"`
	Visibility      string     `json:"visibility"`
	Budget          *float64   `json:"budget"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Builder information
	Builder BuilderInfo `json:"builder"`

	// Application status for this labour user
	HasApplied        bool    `json:"has_applied"`
	ApplicationStatus *string `json:"application_status"` // null if not applied
	ApplicationID     *string `json:"application_id"`     // null if not applied
}

// BuilderInfo represents basic builder information
type BuilderInfo struct {
	BuilderID   string  `json:"builder_id"`
	CompanyName string  `json:"company_name"`
	DisplayName string  `json:"display_name"`
	Location    string  `json:"location"`
	AvatarURL   *string `json:"avatar_url"`
}

// LabourJobsResponse represents the response for labour jobs
type LabourJobsResponse struct {
	Jobs    []LabourJobInfo `json:"jobs"`
	Total   int             `json:"total"`
	Message string          `json:"message"`
}
