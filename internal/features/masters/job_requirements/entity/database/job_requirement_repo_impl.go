package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/job_requirements/models"
	"gorm.io/gorm"
)

// jobRequirementRepository implements JobRequirementRepository
type jobRequirementRepository struct {
	db *gorm.DB
}

// NewJobRequirementRepository creates a new job requirement repository
func NewJobRequirementRepository(db *gorm.DB) JobRequirementRepository {
	return &jobRequirementRepository{
		db: db,
	}
}

// GetAll retrieves all job requirements
func (r *jobRequirementRepository) GetAll(ctx context.Context) ([]*models.JobRequirement, error) {
	var requirements []*models.JobRequirement
	err := r.db.WithContext(ctx).Find(&requirements).Error
	if err != nil {
		return nil, err
	}
	return requirements, nil
}

// GetActive retrieves all active job requirements
func (r *jobRequirementRepository) GetActive(ctx context.Context) ([]*models.JobRequirement, error) {
	var requirements []*models.JobRequirement
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&requirements).Error
	if err != nil {
		return nil, err
	}
	return requirements, nil
}

// GetByID retrieves a job requirement by ID
func (r *jobRequirementRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobRequirement, error) {
	var requirement models.JobRequirement
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&requirement).Error
	if err != nil {
		return nil, err
	}
	return &requirement, nil
}
