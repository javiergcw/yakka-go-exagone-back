package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobsites/models"
)

// JobsiteRepository defines the interface for jobsite database operations
type JobsiteRepository interface {
	Create(ctx context.Context, jobsite *models.Jobsite) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Jobsite, error)
	GetByBuilderID(ctx context.Context, builderID uuid.UUID) ([]*models.Jobsite, error)
	Update(ctx context.Context, jobsite *models.Jobsite) error
	Delete(ctx context.Context, id uuid.UUID) error
}
