package models

import (
	"time"

	"github.com/google/uuid"
)

// Job represents a job posting in the system
type Job struct {
	ID                          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BuilderProfileID            uuid.UUID     `json:"builder_profile_id" gorm:"type:uuid;not null"`
	JobsiteID                   uuid.UUID     `json:"jobsite_id" gorm:"type:uuid;not null"`
	JobTypeID                   uuid.UUID     `json:"job_type_id" gorm:"type:uuid;not null"`
	ManyLabours                 int           `json:"many_labours" gorm:"not null"`
	OngoingWork                 bool          `json:"ongoing_work" gorm:"not null;default:false"`
	WageSiteAllowance           *float64      `json:"wage_site_allowance" gorm:"type:decimal(10,2)"`
	WageLeadingHandAllowance    *float64      `json:"wage_leading_hand_allowance" gorm:"type:decimal(10,2)"`
	WageProductivityAllowance   *float64      `json:"wage_productivity_allowance" gorm:"type:decimal(10,2)"`
	ExtrasOvertimeRate          *float64      `json:"extras_overtime_rate" gorm:"type:decimal(10,2)"`
	WageHourlyRate              *float64      `json:"wage_hourly_rate" gorm:"type:decimal(10,2)"`
	TravelAllowance             *float64      `json:"travel_allowance" gorm:"type:decimal(10,2)"`
	GST                         *float64      `json:"gst" gorm:"type:decimal(10,2)"`
	StartDateWork               *time.Time    `json:"start_date_work" gorm:"type:date"`
	EndDateWork                 *time.Time    `json:"end_date_work" gorm:"type:date"`
	WorkSaturday                bool          `json:"work_saturday" gorm:"not null;default:false"`
	WorkSunday                  bool          `json:"work_sunday" gorm:"not null;default:false"`
	StartTime                   *string       `json:"start_time" gorm:"size:8"` // Format: "HH:MM:SS"
	EndTime                     *string       `json:"end_time" gorm:"size:8"`   // Format: "HH:MM:SS"
	Description                 *string       `json:"description" gorm:"type:text"`
	PaymentDay                  *int          `json:"payment_day" gorm:"type:smallint"` // Day of month for FIXED_DAY
	RequiresSupervisorSignature bool          `json:"requires_supervisor_signature" gorm:"not null;default:false"`
	SupervisorName              *string       `json:"supervisor_name" gorm:"size:100"`
	Visibility                  JobVisibility `json:"visibility" gorm:"type:varchar(20);not null;default:'DRAFT'"`
	PaymentType                 PaymentType   `json:"payment_type" gorm:"type:varchar(20);not null;default:'WEEKLY'"`
	CreatedAt                   time.Time     `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt                   time.Time     `json:"updated_at" gorm:"not null;type:timestamptz"`

	// Relations - loaded separately to avoid circular imports
	// BuilderProfile *BuilderProfile `json:"builder_profile,omitempty" gorm:"foreignKey:BuilderProfileID"`
	// Jobsite        *Jobsite        `json:"jobsite,omitempty" gorm:"foreignKey:JobsiteID"`
	// JobType        *JobType        `json:"job_type,omitempty" gorm:"foreignKey:JobTypeID"`
	JobLicenses []JobLicense `json:"job_licenses,omitempty" gorm:"foreignKey:JobID"`
	JobSkills   []JobSkill   `json:"job_skills,omitempty" gorm:"foreignKey:JobID"`
}

// TableName returns the table name for the Job model
func (Job) TableName() string {
	return "jobs"
}

// JobLicense represents the many-to-many relationship between jobs and licenses
type JobLicense struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobID     uuid.UUID `json:"job_id" gorm:"type:uuid;not null"`
	LicenseID uuid.UUID `json:"license_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;type:timestamptz"`

	// Relations - loaded separately to avoid circular imports
	// Job     *Job     `json:"job,omitempty" gorm:"foreignKey:JobID"`
	// License *License `json:"license,omitempty" gorm:"foreignKey:LicenseID"`
}

// TableName returns the table name for the JobLicense model
func (JobLicense) TableName() string {
	return "job_licenses"
}

// JobSkill represents the many-to-many relationship between jobs and skills
type JobSkill struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobID              uuid.UUID  `json:"job_id" gorm:"type:uuid;not null"`
	SkillCategoryID    *uuid.UUID `json:"skill_category_id" gorm:"type:uuid"`
	SkillSubcategoryID *uuid.UUID `json:"skill_subcategory_id" gorm:"type:uuid"`
	CreatedAt          time.Time  `json:"created_at" gorm:"not null;type:timestamptz"`

	// Relations - loaded separately to avoid circular imports
	// Job              *Job              `json:"job,omitempty" gorm:"foreignKey:JobID"`
	// SkillCategory    *SkillCategory    `json:"skill_category,omitempty" gorm:"foreignKey:SkillCategoryID"`
	// SkillSubcategory *SkillSubcategory `json:"skill_subcategory,omitempty" gorm:"foreignKey:SkillSubcategoryID"`
}

// TableName returns the table name for the JobSkill model
func (JobSkill) TableName() string {
	return "job_skills"
}
