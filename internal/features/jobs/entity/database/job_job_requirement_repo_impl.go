package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
	"gorm.io/gorm"
)

// jobJobRequirementRepository implements JobJobRequirementRepository
type jobJobRequirementRepository struct {
	db *gorm.DB
}

// NewJobJobRequirementRepository creates a new job job requirement repository
func NewJobJobRequirementRepository(db *gorm.DB) JobJobRequirementRepository {
	return &jobJobRequirementRepository{
		db: db,
	}
}

// Create creates a new job job requirement relationship
func (r *jobJobRequirementRepository) Create(ctx context.Context, jobJobRequirement *models.JobJobRequirement) error {
	return r.db.WithContext(ctx).Create(jobJobRequirement).Error
}

// GetByJobID retrieves all job job requirements for a specific job
func (r *jobJobRequirementRepository) GetByJobID(ctx context.Context, jobID uuid.UUID) ([]*models.JobJobRequirement, error) {
	var jobJobRequirements []*models.JobJobRequirement
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Find(&jobJobRequirements).Error
	if err != nil {
		return nil, err
	}
	return jobJobRequirements, nil
}

// DeleteByJobID deletes all job job requirements for a specific job
func (r *jobJobRequirementRepository) DeleteByJobID(ctx context.Context, jobID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ?", jobID).Delete(&models.JobJobRequirement{}).Error
}

// DeleteByJobAndRequirement deletes a specific job job requirement relationship
func (r *jobJobRequirementRepository) DeleteByJobAndRequirement(ctx context.Context, jobID, jobRequirementID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ? AND job_requirement_id = ?", jobID, jobRequirementID).Delete(&models.JobJobRequirement{}).Error
}
