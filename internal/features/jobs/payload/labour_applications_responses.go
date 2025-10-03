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

// LabourApplicationInfo represents detailed application information for a labour user
type LabourApplicationInfo struct {
	ApplicationID string    `json:"application_id"`
	JobID         string    `json:"job_id"`
	Status        string    `json:"status"`
	CoverLetter   *string   `json:"cover_letter"`
	ExpectedRate  *float64  `json:"expected_rate"`
	ResumeURL     *string   `json:"resume_url"`
	AppliedAt     time.Time `json:"applied_at"`
	Job           JobInfo   `json:"job"`
}

// JobInfo represents basic job information for labour applications
type JobInfo struct {
	ID             string                  `json:"id"`
	Description    string                  `json:"description"`
	StartDate      *time.Time              `json:"start_date"`
	EndDate        *time.Time              `json:"end_date"`
	WageHourlyRate *float64                `json:"wage_hourly_rate"`
	Visibility     string                  `json:"visibility"`
	CreatedAt      time.Time               `json:"created_at"`
	BuilderProfile *BuilderProfileInfo     `json:"builder_profile,omitempty"`
	Jobsite        *JobsiteApplicationInfo `json:"jobsite,omitempty"`
	JobType        *JobTypeInfo            `json:"job_type,omitempty"`
}

// BuilderProfileInfo represents basic builder profile information
type BuilderProfileInfo struct {
	ID          string  `json:"id"`
	CompanyName *string `json:"company_name"`
	DisplayName string  `json:"display_name"`
	Location    *string `json:"location"`
}

// JobsiteApplicationInfo represents basic jobsite information for applications
type JobsiteApplicationInfo struct {
	ID          string  `json:"id"`
	Name        *string `json:"name"`
	Address     string  `json:"address"`
	City        *string `json:"city"`
	Suburb      *string `json:"suburb"`
	Description *string `json:"description"`
}

// JobTypeInfo represents basic job type information
type JobTypeInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// LabourApplicantsResponse represents the response for labour applicants
type LabourApplicantsResponse struct {
	Applications []LabourApplicationInfo `json:"applications"`
	Total        int                     `json:"total"`
	Message      string                  `json:"message"`
}
