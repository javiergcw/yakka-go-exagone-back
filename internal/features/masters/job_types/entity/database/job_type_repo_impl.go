package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/job_types/models"
	"gorm.io/gorm"
)

// jobTypeRepository implements JobTypeRepository
type jobTypeRepository struct {
	db *gorm.DB
}

// NewJobTypeRepository creates a new job type repository
func NewJobTypeRepository(db *gorm.DB) JobTypeRepository {
	return &jobTypeRepository{
		db: db,
	}
}

// GetAll retrieves all job types
func (r *jobTypeRepository) GetAll(ctx context.Context) ([]*models.JobType, error) {
	var types []*models.JobType
	err := r.db.WithContext(ctx).Find(&types).Error
	if err != nil {
		return nil, err
	}
	return types, nil
}

// GetActive retrieves all active job types
func (r *jobTypeRepository) GetActive(ctx context.Context) ([]*models.JobType, error) {
	var types []*models.JobType
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&types).Error
	if err != nil {
		return nil, err
	}
	return types, nil
}

// GetByID retrieves a job type by ID
func (r *jobTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobType, error) {
	var jobType models.JobType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&jobType).Error
	if err != nil {
		return nil, err
	}
	return &jobType, nil
}
