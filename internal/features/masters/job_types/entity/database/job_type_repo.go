package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/job_types/models"
)

// JobTypeRepository defines the interface for job type operations
type JobTypeRepository interface {
	GetAll(ctx context.Context) ([]*models.JobType, error)
	GetActive(ctx context.Context) ([]*models.JobType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.JobType, error)
}
