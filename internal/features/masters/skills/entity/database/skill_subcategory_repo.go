package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/skills/models"
)

type SkillSubcategoryRepository interface {
	Create(ctx context.Context, subcategory *models.SkillSubcategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SkillSubcategory, error)
	GetAll(ctx context.Context) ([]*models.SkillSubcategory, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*models.SkillSubcategory, error)
	Update(ctx context.Context, subcategory *models.SkillSubcategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}
