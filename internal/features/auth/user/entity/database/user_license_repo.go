package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user/models"
)

// UserLicenseRepository defines the interface for user license database operations
type UserLicenseRepository interface {
	Create(ctx context.Context, license *models.UserLicense) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.UserLicense, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserLicense, error)
	Update(ctx context.Context, license *models.UserLicense) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	CreateBatch(ctx context.Context, licenses []*models.UserLicense) error
}
