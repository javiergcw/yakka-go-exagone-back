package payload

import "time"

// CreateJobAssignmentRequest represents the request to create a job assignment
type CreateJobAssignmentRequest struct {
	JobID         string     `json:"job_id" validate:"required,uuid"`
	LabourUserID  string     `json:"labour_user_id" validate:"required,uuid"`
	ApplicationID string     `json:"application_id" validate:"required,uuid"`
	StartDate     *time.Time `json:"start_date" validate:"omitempty"`
	EndDate       *time.Time `json:"end_date" validate:"omitempty"`
}

// UpdateJobAssignmentRequest represents the request to update a job assignment
type UpdateJobAssignmentRequest struct {
	StartDate *time.Time `json:"start_date" validate:"omitempty"`
	EndDate   *time.Time `json:"end_date" validate:"omitempty"`
	Status    *string    `json:"status" validate:"omitempty,oneof=ACTIVE COMPLETED CANCELLED"`
}

// GetJobAssignmentsRequest represents the request to get job assignments with filters
type GetJobAssignmentsRequest struct {
	JobID         *string `json:"job_id" form:"job_id"`
	LabourUserID  *string `json:"labour_user_id" form:"labour_user_id"`
	ApplicationID *string `json:"application_id" form:"application_id"`
	Status        *string `json:"status" form:"status"`
	Page          int     `json:"page" form:"page" validate:"min=1"`
	Limit         int     `json:"limit" form:"limit" validate:"min=1,max=100"`
}

// CompleteAssignmentRequest represents the request to complete an assignment
type CompleteAssignmentRequest struct {
	EndDate *time.Time `json:"end_date" validate:"omitempty"`
}

// CancelAssignmentRequest represents the request to cancel an assignment
type CancelAssignmentRequest struct {
	Reason *string `json:"reason"`
}
