package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/labour_profiles/models"
	"gorm.io/gorm"
)

type labourProfileSkillRepository struct {
	db *gorm.DB
}

// NewLabourProfileSkillRepository creates a new labour profile skill repository
func NewLabourProfileSkillRepository(db *gorm.DB) LabourProfileSkillRepository {
	return &labourProfileSkillRepository{db: db}
}

func (r *labourProfileSkillRepository) Create(ctx context.Context, skill *models.LabourProfileSkill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *labourProfileSkillRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LabourProfileSkill, error) {
	var skill models.LabourProfileSkill
	err := r.db.WithContext(ctx).Preload("Category").Preload("Subcategory").Preload("ExperienceLevel").First(&skill, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *labourProfileSkillRepository) GetByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) ([]*models.LabourProfileSkill, error) {
	var skills []*models.LabourProfileSkill
	err := r.db.WithContext(ctx).Preload("Category").Preload("Subcategory").Preload("ExperienceLevel").Where("labour_profile_id = ?", labourProfileID).Find(&skills).Error
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *labourProfileSkillRepository) Update(ctx context.Context, skill *models.LabourProfileSkill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

func (r *labourProfileSkillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.LabourProfileSkill{}, "id = ?", id).Error
}

func (r *labourProfileSkillRepository) DeleteByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.LabourProfileSkill{}, "labour_profile_id = ?", labourProfileID).Error
}

func (r *labourProfileSkillRepository) CreateBatch(ctx context.Context, skills []*models.LabourProfileSkill) error {
	return r.db.WithContext(ctx).CreateInBatches(skills, 100).Error
}
