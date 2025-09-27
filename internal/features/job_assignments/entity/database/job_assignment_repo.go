package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_assignments/models"
)

// JobAssignmentRepository defines the interface for job assignment data operations
type JobAssignmentRepository interface {
	// Create creates a new job assignment
	Create(ctx context.Context, assignment *models.JobAssignment) error

	// GetByID retrieves a job assignment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.JobAssignment, error)

	// GetByApplicationID retrieves a job assignment by application ID
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) (*models.JobAssignment, error)

	// Update updates an existing job assignment
	Update(ctx context.Context, assignment *models.JobAssignment) error

	// Delete deletes a job assignment
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByJobID retrieves all assignments for a specific job
	GetByJobID(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error)

	// GetByLabourUserID retrieves all assignments for a specific labour user
	GetByLabourUserID(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error)

	// GetByStatus retrieves assignments by status
	GetByStatus(ctx context.Context, status models.AssignmentStatus, page, limit int) ([]*models.JobAssignment, int64, error)

	// GetWithFilters retrieves assignments with multiple filters
	GetWithFilters(ctx context.Context, jobID, labourUserID, applicationID *uuid.UUID, status *models.AssignmentStatus, page, limit int) ([]*models.JobAssignment, int64, error)

	// UpdateStatus updates the status of a job assignment
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.AssignmentStatus) error

	// CompleteAssignment completes an assignment
	CompleteAssignment(ctx context.Context, id uuid.UUID, endDate *time.Time) error

	// CancelAssignment cancels an assignment
	CancelAssignment(ctx context.Context, id uuid.UUID) error

	// CheckAssignmentExists checks if an assignment already exists for an application
	CheckAssignmentExists(ctx context.Context, applicationID uuid.UUID) (bool, error)
}
