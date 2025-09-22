package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/experience_levels/models"
)

type ExperienceLevelRepository interface {
	Create(ctx context.Context, experienceLevel *models.ExperienceLevel) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ExperienceLevel, error)
	GetAll(ctx context.Context) ([]*models.ExperienceLevel, error)
	Update(ctx context.Context, experienceLevel *models.ExperienceLevel) error
	Delete(ctx context.Context, id uuid.UUID) error
}
