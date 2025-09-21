package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/builder_profiles/models"
)

// BuilderProfileRepository defines the interface for builder profile database operations
type BuilderProfileRepository interface {
	Create(ctx context.Context, profile *models.BuilderProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BuilderProfile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.BuilderProfile, error)
	Update(ctx context.Context, profile *models.BuilderProfile) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
