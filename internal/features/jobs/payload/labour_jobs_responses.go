package payload

import "time"

// LabourJobInfo represents a job with application status for a labour user
type LabourJobInfo struct {
	JobID           string     `json:"job_id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Location        string     `json:"location"`
	JobType         string     `json:"job_type"`
	ManyLabours     int        `json:"many_labours"`
	ExperienceLevel string     `json:"experience_level"`
	Status          string     `json:"status"`
	Visibility      string     `json:"visibility"`
	TotalWage       *float64   `json:"total_wage"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Builder information
	Builder BuilderInfo `json:"builder"`

	// Jobsite information
	Jobsite *JobsiteInfo `json:"jobsite,omitempty"`

	// Skills information
	Skills []JobSkillInfo `json:"skills,omitempty"`

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

// JobsiteInfo represents jobsite information for labour jobs
type JobsiteInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        *string `json:"city,omitempty"`
	Suburb      *string `json:"suburb,omitempty"`
	Description *string `json:"description,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Phone       *string `json:"phone,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// JobSkillInfo represents skill information for labour jobs
type JobSkillInfo struct {
	ID                 string                `json:"id"`
	JobID              string                `json:"job_id"`
	SkillCategoryID    *string               `json:"skill_category_id,omitempty"`
	SkillSubcategoryID *string               `json:"skill_subcategory_id,omitempty"`
	SkillCategory      *SkillCategoryInfo    `json:"skill_category,omitempty"`
	SkillSubcategory   *SkillSubcategoryInfo `json:"skill_subcategory,omitempty"`
	CreatedAt          string                `json:"created_at"`
}

// SkillCategoryInfo represents skill category information
type SkillCategoryInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// SkillSubcategoryInfo represents skill subcategory information
type SkillSubcategoryInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// LabourJobsResponse represents the response for labour jobs
type LabourJobsResponse struct {
	Jobs    []LabourJobInfo `json:"jobs"`
	Total   int             `json:"total"`
	Message string          `json:"message"`
}
