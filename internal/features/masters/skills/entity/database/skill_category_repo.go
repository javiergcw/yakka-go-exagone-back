package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/skills/models"
)

type SkillCategoryRepository interface {
	Create(ctx context.Context, category *models.SkillCategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SkillCategory, error)
	GetAll(ctx context.Context) ([]*models.SkillCategory, error)
	Update(ctx context.Context, category *models.SkillCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}
