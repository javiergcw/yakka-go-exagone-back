package payload

import "time"

// BuilderApplicantDecisionRequest represents the request to hire or reject an applicant
type BuilderApplicantDecisionRequest struct {
	ApplicationID string     `json:"application_id" validate:"required,uuid"`
	Hired         *bool      `json:"hired" validate:"required"`
	StartDate     *time.Time `json:"start_date" validate:"omitempty"`
	EndDate       *time.Time `json:"end_date" validate:"omitempty"`
	Reason        *string    `json:"reason" validate:"omitempty"`
}
