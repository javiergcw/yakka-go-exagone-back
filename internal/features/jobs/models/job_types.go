package models

// JobVisibility represents the visibility status of a job
type JobVisibility string

const (
	JobVisibilityPublic   JobVisibility = "PUBLIC"
	JobVisibilityPrivate  JobVisibility = "PRIVATE"
	JobVisibilityBanned   JobVisibility = "BANNED"
	JobVisibilityArchived JobVisibility = "ARCHIVED"
	JobVisibilityDraft    JobVisibility = "DRAFT"
)

// PaymentType represents the payment frequency for a job
type PaymentType string

const (
	PaymentTypeFixedDay    PaymentType = "FIXED_DAY"
	PaymentTypeWeekly      PaymentType = "WEEKLY"
	PaymentTypeFortnightly PaymentType = "FORTNIGHTLY"
)
