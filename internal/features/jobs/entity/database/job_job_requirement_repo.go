package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
)

// JobJobRequirementRepository defines the interface for job job requirement operations
type JobJobRequirementRepository interface {
	Create(ctx context.Context, jobJobRequirement *models.JobJobRequirement) error
	GetByJobID(ctx context.Context, jobID uuid.UUID) ([]*models.JobJobRequirement, error)
	DeleteByJobID(ctx context.Context, jobID uuid.UUID) error
	DeleteByJobAndRequirement(ctx context.Context, jobID, jobRequirementID uuid.UUID) error
}
