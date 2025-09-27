package models

// ApplicationStatus represents the status of a job application
type ApplicationStatus string

const (
	ApplicationStatusApplied   ApplicationStatus = "APPLIED"
	ApplicationStatusReviewed  ApplicationStatus = "REVIEWED"
	ApplicationStatusAccepted  ApplicationStatus = "ACCEPTED"
	ApplicationStatusRejected  ApplicationStatus = "REJECTED"
	ApplicationStatusWithdrawn ApplicationStatus = "WITHDRAWN"
)

// IsValid checks if the application status is valid
func (s ApplicationStatus) IsValid() bool {
	switch s {
	case ApplicationStatusApplied, ApplicationStatusReviewed, ApplicationStatusAccepted, ApplicationStatusRejected, ApplicationStatusWithdrawn:
		return true
	default:
		return false
	}
}

// String returns the string representation of the status
func (s ApplicationStatus) String() string {
	return string(s)
}
