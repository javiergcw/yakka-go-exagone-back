package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/experience_levels/models"
	"gorm.io/gorm"
)

type experienceLevelRepository struct {
	db *gorm.DB
}

func NewExperienceLevelRepository(db *gorm.DB) ExperienceLevelRepository {
	return &experienceLevelRepository{db: db}
}

func (r *experienceLevelRepository) Create(ctx context.Context, experienceLevel *models.ExperienceLevel) error {
	return r.db.WithContext(ctx).Create(experienceLevel).Error
}

func (r *experienceLevelRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ExperienceLevel, error) {
	var experienceLevel models.ExperienceLevel
	err := r.db.WithContext(ctx).First(&experienceLevel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &experienceLevel, nil
}

func (r *experienceLevelRepository) GetAll(ctx context.Context) ([]*models.ExperienceLevel, error) {
	var experienceLevels []*models.ExperienceLevel
	err := r.db.WithContext(ctx).Find(&experienceLevels).Error
	if err != nil {
		return nil, err
	}
	return experienceLevels, nil
}

func (r *experienceLevelRepository) Update(ctx context.Context, experienceLevel *models.ExperienceLevel) error {
	return r.db.WithContext(ctx).Save(experienceLevel).Error
}

func (r *experienceLevelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ExperienceLevel{}, "id = ?", id).Error
}
