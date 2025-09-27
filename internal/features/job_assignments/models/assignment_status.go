package models

// AssignmentStatus represents the status of a job assignment
type AssignmentStatus string

const (
	AssignmentStatusActive    AssignmentStatus = "ACTIVE"
	AssignmentStatusCompleted AssignmentStatus = "COMPLETED"
	AssignmentStatusCancelled AssignmentStatus = "CANCELLED"
)

// IsValid checks if the assignment status is valid
func (s AssignmentStatus) IsValid() bool {
	switch s {
	case AssignmentStatusActive, AssignmentStatusCompleted, AssignmentStatusCancelled:
		return true
	default:
		return false
	}
}

// String returns the string representation of the status
func (s AssignmentStatus) String() string {
	return string(s)
}
