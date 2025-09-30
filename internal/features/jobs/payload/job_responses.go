package payload

import (
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
)

// JobResponse represents a job in API responses
type JobResponse struct {
	ID                          uuid.UUID            `json:"id"`
	ManyLabours                 int                  `json:"many_labours"`
	OngoingWork                 bool                 `json:"ongoing_work"`
	WageSiteAllowance           *float64             `json:"wage_site_allowance"`
	WageLeadingHandAllowance    *float64             `json:"wage_leading_hand_allowance"`
	WageProductivityAllowance   *float64             `json:"wage_productivity_allowance"`
	ExtrasOvertimeRate          *float64             `json:"extras_overtime_rate"`
	WageHourlyRate              *float64             `json:"wage_hourly_rate"`
	TravelAllowance             *float64             `json:"travel_allowance"`
	GST                         *float64             `json:"gst"`
	StartDateWork               *time.Time           `json:"start_date_work"`
	EndDateWork                 *time.Time           `json:"end_date_work"`
	WorkSaturday                bool                 `json:"work_saturday"`
	WorkSunday                  bool                 `json:"work_sunday"`
	StartTime                   *string              `json:"start_time"`
	EndTime                     *string              `json:"end_time"`
	Description                 *string              `json:"description"`
	PaymentDay                  *int                 `json:"payment_day"`
	RequiresSupervisorSignature bool                 `json:"requires_supervisor_signature"`
	SupervisorName              *string              `json:"supervisor_name"`
	Visibility                  models.JobVisibility `json:"visibility"`
	PaymentType                 models.PaymentType   `json:"payment_type"`
	CreatedAt                   time.Time            `json:"created_at"`
	UpdatedAt                   time.Time            `json:"updated_at"`

	// Relations
	BuilderProfile  *BuilderProfileResponse  `json:"builder_profile,omitempty"`
	Jobsite         *JobsiteResponse         `json:"jobsite"`
	JobType         *JobTypeResponse         `json:"job_type"`
	JobLicenses     []JobLicenseResponse     `json:"job_licenses"`
	JobSkills       []JobSkillResponse       `json:"job_skills"`
	JobRequirements []JobRequirementResponse `json:"job_requirements"`
}

// BuilderProfileResponse represents a builder profile in responses
type BuilderProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	CompanyName string    `json:"company_name"`
	DisplayName *string   `json:"display_name,omitempty"`
	Location    *string   `json:"location,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	Email       *string   `json:"email,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JobsiteResponse represents a jobsite in responses
type JobsiteResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        *string   `json:"city,omitempty"`
	Suburb      *string   `json:"suburb,omitempty"`
	Description *string   `json:"description,omitempty"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Phone       *string   `json:"phone,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JobTypeResponse represents a job type in responses
type JobTypeResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JobLicenseResponse represents a job license in responses
type JobLicenseResponse struct {
	ID        uuid.UUID        `json:"id"`
	JobID     uuid.UUID        `json:"job_id"`
	LicenseID uuid.UUID        `json:"license_id"`
	License   *LicenseResponse `json:"license,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

// LicenseResponse represents a license in responses
type LicenseResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

// JobSkillResponse represents a job skill in responses
type JobSkillResponse struct {
	ID                 uuid.UUID                 `json:"id"`
	JobID              uuid.UUID                 `json:"job_id"`
	SkillCategoryID    *uuid.UUID                `json:"skill_category_id"`
	SkillSubcategoryID *uuid.UUID                `json:"skill_subcategory_id"`
	SkillCategory      *SkillCategoryResponse    `json:"skill_category,omitempty"`
	SkillSubcategory   *SkillSubcategoryResponse `json:"skill_subcategory,omitempty"`
	CreatedAt          time.Time                 `json:"created_at"`
}

// SkillCategoryResponse represents a skill category in responses
type SkillCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

// SkillSubcategoryResponse represents a skill subcategory in responses
type SkillSubcategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

// CreateJobResponse represents the response when creating a job
type CreateJobResponse struct {
	Job     JobResponse `json:"job"`
	Message string      `json:"message"`
}

// GetJobResponse represents the response when getting a job
type GetJobResponse struct {
	Job     JobResponse `json:"job"`
	Message string      `json:"message"`
}

// GetJobsResponse represents the response when getting multiple jobs
type GetJobsResponse struct {
	Jobs    []JobResponse `json:"jobs"`
	Message string        `json:"message"`
}

// UpdateJobResponse represents the response when updating a job
type UpdateJobResponse struct {
	Job     JobResponse `json:"job"`
	Message string      `json:"message"`
}

// DeleteJobResponse represents the response when deleting a job
type DeleteJobResponse struct {
	Message string `json:"message"`
}

// UpdateJobVisibilityResponse represents the response when updating job visibility
type UpdateJobVisibilityResponse struct {
	Job     JobResponse `json:"job"`
	Message string      `json:"message"`
}

// LabourJobDetailResponse represents the response for labour job detail with application info
type LabourJobDetailResponse struct {
	Job         JobResponse         `json:"job"`
	Application *JobApplicationInfo `json:"application"`
	Message     string              `json:"message"`
}

// JobApplicationInfo represents application information for a labour user
type JobApplicationInfo struct {
	ID           uuid.UUID  `json:"id"`
	Status       string     `json:"status"`
	CoverLetter  *string    `json:"cover_letter,omitempty"`
	ExpectedRate *float64   `json:"expected_rate,omitempty"`
	ResumeURL    *string    `json:"resume_url,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	WithdrawnAt  *time.Time `json:"withdrawn_at,omitempty"`
}

// JobRequirementResponse represents a job requirement in responses
type JobRequirementResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
