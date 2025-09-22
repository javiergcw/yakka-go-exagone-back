package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/licenses/models"
)

type LicenseRepository interface {
	Create(ctx context.Context, license *models.License) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.License, error)
	GetAll(ctx context.Context) ([]*models.License, error)
	Update(ctx context.Context, license *models.License) error
	Delete(ctx context.Context, id uuid.UUID) error
}
