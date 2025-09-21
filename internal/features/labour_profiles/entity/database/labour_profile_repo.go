package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/labour_profiles/models"
)

// LabourProfileRepository defines the interface for labour profile database operations
type LabourProfileRepository interface {
	Create(ctx context.Context, profile *models.LabourProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LabourProfile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.LabourProfile, error)
	Update(ctx context.Context, profile *models.LabourProfile) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
