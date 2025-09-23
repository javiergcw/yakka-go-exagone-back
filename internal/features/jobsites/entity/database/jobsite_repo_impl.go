package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobsites/models"
	"gorm.io/gorm"
)

// JobsiteRepositoryImpl implements the JobsiteRepository interface
type JobsiteRepositoryImpl struct {
	db *gorm.DB
}

// NewJobsiteRepository creates a new instance of JobsiteRepositoryImpl
func NewJobsiteRepositoryImpl(db *gorm.DB) JobsiteRepository {
	return &JobsiteRepositoryImpl{db: db}
}

// Create creates a new jobsite
func (r *JobsiteRepositoryImpl) Create(ctx context.Context, jobsite *models.Jobsite) error {
	return r.db.WithContext(ctx).Create(jobsite).Error
}

// GetByID retrieves a jobsite by its ID
func (r *JobsiteRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Jobsite, error) {
	var jobsite models.Jobsite
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&jobsite).Error
	if err != nil {
		return nil, err
	}
	return &jobsite, nil
}

// GetByBuilderID retrieves all jobsites for a specific builder
func (r *JobsiteRepositoryImpl) GetByBuilderID(ctx context.Context, builderID uuid.UUID) ([]*models.Jobsite, error) {
	var jobsites []*models.Jobsite
	err := r.db.WithContext(ctx).Where("builder_id = ?", builderID).Find(&jobsites).Error
	if err != nil {
		return nil, err
	}
	return jobsites, nil
}

// Update updates an existing jobsite
func (r *JobsiteRepositoryImpl) Update(ctx context.Context, jobsite *models.Jobsite) error {
	return r.db.WithContext(ctx).Save(jobsite).Error
}

// Delete deletes a jobsite by its ID
func (r *JobsiteRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Jobsite{}).Error
}
