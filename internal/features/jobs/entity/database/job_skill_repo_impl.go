package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
	"gorm.io/gorm"
)

// jobSkillRepository implements JobSkillRepository
type jobSkillRepository struct {
	db *gorm.DB
}

// NewJobSkillRepository creates a new job skill repository
func NewJobSkillRepository(db *gorm.DB) JobSkillRepository {
	return &jobSkillRepository{
		db: db,
	}
}

// Create creates a new job skill relationship
func (r *jobSkillRepository) Create(ctx context.Context, jobSkill *models.JobSkill) error {
	return r.db.WithContext(ctx).Create(jobSkill).Error
}

// GetByJobID retrieves all job skills for a specific job
func (r *jobSkillRepository) GetByJobID(ctx context.Context, jobID uuid.UUID) ([]*models.JobSkill, error) {
	var jobSkills []*models.JobSkill
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Find(&jobSkills).Error
	if err != nil {
		return nil, err
	}
	return jobSkills, nil
}

// DeleteByJobID deletes all job skills for a specific job
func (r *jobSkillRepository) DeleteByJobID(ctx context.Context, jobID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ?", jobID).Delete(&models.JobSkill{}).Error
}

// DeleteByJobAndSkill deletes a specific job skill relationship
func (r *jobSkillRepository) DeleteByJobAndSkill(ctx context.Context, jobID, skillCategoryID, skillSubcategoryID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ? AND skill_category_id = ? AND skill_subcategory_id = ?", jobID, skillCategoryID, skillSubcategoryID).Delete(&models.JobSkill{}).Error
}
