package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/skills/models"
	"gorm.io/gorm"
)

type skillCategoryRepository struct {
	db *gorm.DB
}

func NewSkillCategoryRepository(db *gorm.DB) SkillCategoryRepository {
	return &skillCategoryRepository{db: db}
}

func (r *skillCategoryRepository) Create(ctx context.Context, category *models.SkillCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *skillCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SkillCategory, error) {
	var category models.SkillCategory
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *skillCategoryRepository) GetAll(ctx context.Context) ([]*models.SkillCategory, error) {
	var categories []*models.SkillCategory
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *skillCategoryRepository) Update(ctx context.Context, category *models.SkillCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *skillCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.SkillCategory{}, "id = ?", id).Error
}
