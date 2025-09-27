package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_applications/models"
)

// JobApplicationRepository defines the interface for job application data operations
type JobApplicationRepository interface {
	// Create creates a new job application
	Create(ctx context.Context, application *models.JobApplication) error

	// GetByID retrieves a job application by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.JobApplication, error)

	// GetByJobAndLabourUser retrieves a job application by job ID and labour user ID
	GetByJobAndLabourUser(ctx context.Context, jobID, labourUserID uuid.UUID) (*models.JobApplication, error)

	// Update updates an existing job application
	Update(ctx context.Context, application *models.JobApplication) error

	// Delete deletes a job application
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByJobID retrieves all applications for a specific job
	GetByJobID(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error)

	// GetByLabourUserID retrieves all applications for a specific labour user
	GetByLabourUserID(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error)

	// GetByStatus retrieves applications by status
	GetByStatus(ctx context.Context, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error)

	// GetWithFilters retrieves applications with multiple filters
	GetWithFilters(ctx context.Context, jobID, labourUserID *uuid.UUID, status *models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error)

	// UpdateStatus updates the status of a job application
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.ApplicationStatus) error

	// WithdrawApplication withdraws an application
	WithdrawApplication(ctx context.Context, id uuid.UUID) error

	// CheckApplicationExists checks if an application already exists for a job and user
	CheckApplicationExists(ctx context.Context, jobID, labourUserID uuid.UUID) (bool, error)
}
