package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
)

// JobRepository defines the interface for job operations
type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Job, error)
	GetByBuilderProfileID(ctx context.Context, builderProfileID uuid.UUID) ([]*models.Job, error)
	GetByJobsiteID(ctx context.Context, jobsiteID uuid.UUID) ([]*models.Job, error)
	GetByVisibility(ctx context.Context, visibility models.JobVisibility) ([]*models.Job, error)
	GetAll(ctx context.Context) ([]*models.Job, error)
	Update(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithRelations(ctx context.Context, id uuid.UUID) (*models.Job, error)
}
