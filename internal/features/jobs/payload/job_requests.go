package payload

import (
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
)

// CreateJobRequest represents the request to create a new job
type CreateJobRequest struct {
	BuilderProfileID            uuid.UUID            `json:"-"` // Set internally from JWT
	JobsiteID                   uuid.UUID            `json:"jobsite_id" validate:"required"`
	JobTypeID                   uuid.UUID            `json:"job_type_id" validate:"required"`
	ManyLabours                 int                  `json:"many_labours" validate:"required,min=1"`
	OngoingWork                 bool                 `json:"ongoing_work"`
	WageSiteAllowance           *float64             `json:"wage_site_allowance"`
	WageLeadingHandAllowance    *float64             `json:"wage_leading_hand_allowance"`
	WageProductivityAllowance   *float64             `json:"wage_productivity_allowance"`
	ExtrasOvertimeRate          *float64             `json:"extras_overtime_rate"`
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
	LicenseIDs                  []uuid.UUID          `json:"license_ids"`
	SkillCategoryIDs            []uuid.UUID          `json:"skill_category_ids"`
	SkillSubcategoryIDs         []uuid.UUID          `json:"skill_subcategory_ids"`
}

// UpdateJobVisibilityRequest represents the request to update job visibility
type UpdateJobVisibilityRequest struct {
	Visibility models.JobVisibility `json:"visibility" validate:"required,oneof=DRAFT PUBLIC PRIVATE"`
}

// UpdateJobRequest represents the request to update a job
type UpdateJobRequest struct {
	ManyLabours                 *int                  `json:"many_labours"`
	OngoingWork                 *bool                 `json:"ongoing_work"`
	WageSiteAllowance           *float64              `json:"wage_site_allowance"`
	WageLeadingHandAllowance    *float64              `json:"wage_leading_hand_allowance"`
	WageProductivityAllowance   *float64              `json:"wage_productivity_allowance"`
	ExtrasOvertimeRate          *float64              `json:"extras_overtime_rate"`
	StartDateWork               *time.Time            `json:"start_date_work"`
	EndDateWork                 *time.Time            `json:"end_date_work"`
	WorkSaturday                *bool                 `json:"work_saturday"`
	WorkSunday                  *bool                 `json:"work_sunday"`
	StartTime                   *string               `json:"start_time"`
	EndTime                     *string               `json:"end_time"`
	Description                 *string               `json:"description"`
	PaymentDay                  *int                  `json:"payment_day"`
	RequiresSupervisorSignature *bool                 `json:"requires_supervisor_signature"`
	SupervisorName              *string               `json:"supervisor_name"`
	Visibility                  *models.JobVisibility `json:"visibility"`
	PaymentType                 *models.PaymentType   `json:"payment_type"`
	LicenseIDs                  []uuid.UUID           `json:"license_ids"`
	SkillCategoryIDs            []uuid.UUID           `json:"skill_category_ids"`
	SkillSubcategoryIDs         []uuid.UUID           `json:"skill_subcategory_ids"`
}

// GetJobsByBuilderRequest represents the request to get jobs by builder
type GetJobsByBuilderRequest struct {
	BuilderProfileID uuid.UUID             `json:"builder_profile_id" validate:"required"`
	Visibility       *models.JobVisibility `json:"visibility"`
}

// GetJobsByJobsiteRequest represents the request to get jobs by jobsite
type GetJobsByJobsiteRequest struct {
	JobsiteID  uuid.UUID             `json:"jobsite_id" validate:"required"`
	Visibility *models.JobVisibility `json:"visibility"`
}
