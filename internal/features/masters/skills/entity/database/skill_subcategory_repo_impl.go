package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/skills/models"
	"gorm.io/gorm"
)

type skillSubcategoryRepository struct {
	db *gorm.DB
}

func NewSkillSubcategoryRepository(db *gorm.DB) SkillSubcategoryRepository {
	return &skillSubcategoryRepository{db: db}
}

func (r *skillSubcategoryRepository) Create(ctx context.Context, subcategory *models.SkillSubcategory) error {
	return r.db.WithContext(ctx).Create(subcategory).Error
}

func (r *skillSubcategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SkillSubcategory, error) {
	var subcategory models.SkillSubcategory
	err := r.db.WithContext(ctx).Preload("Category").First(&subcategory, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

func (r *skillSubcategoryRepository) GetAll(ctx context.Context) ([]*models.SkillSubcategory, error) {
	var subcategories []*models.SkillSubcategory
	err := r.db.WithContext(ctx).Preload("Category").Find(&subcategories).Error
	if err != nil {
		return nil, err
	}
	return subcategories, nil
}

func (r *skillSubcategoryRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*models.SkillSubcategory, error) {
	var subcategories []*models.SkillSubcategory
	err := r.db.WithContext(ctx).Preload("Category").Where("category_id = ?", categoryID).Find(&subcategories).Error
	if err != nil {
		return nil, err
	}
	return subcategories, nil
}

func (r *skillSubcategoryRepository) Update(ctx context.Context, subcategory *models.SkillSubcategory) error {
	return r.db.WithContext(ctx).Save(subcategory).Error
}

func (r *skillSubcategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.SkillSubcategory{}, "id = ?", id).Error
}
