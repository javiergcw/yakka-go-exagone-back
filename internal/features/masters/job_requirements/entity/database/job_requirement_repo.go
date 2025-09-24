package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/job_requirements/models"
)

// JobRequirementRepository defines the interface for job requirement operations
type JobRequirementRepository interface {
	GetAll(ctx context.Context) ([]*models.JobRequirement, error)
	GetActive(ctx context.Context) ([]*models.JobRequirement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.JobRequirement, error)
}
