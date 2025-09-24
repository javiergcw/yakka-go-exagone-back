package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
	"gorm.io/gorm"
)

// jobLicenseRepository implements JobLicenseRepository
type jobLicenseRepository struct {
	db *gorm.DB
}

// NewJobLicenseRepository creates a new job license repository
func NewJobLicenseRepository(db *gorm.DB) JobLicenseRepository {
	return &jobLicenseRepository{
		db: db,
	}
}

// Create creates a new job license relationship
func (r *jobLicenseRepository) Create(ctx context.Context, jobLicense *models.JobLicense) error {
	return r.db.WithContext(ctx).Create(jobLicense).Error
}

// GetByJobID retrieves all job licenses for a specific job
func (r *jobLicenseRepository) GetByJobID(ctx context.Context, jobID uuid.UUID) ([]*models.JobLicense, error) {
	var jobLicenses []*models.JobLicense
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Find(&jobLicenses).Error
	if err != nil {
		return nil, err
	}
	return jobLicenses, nil
}

// DeleteByJobID deletes all job licenses for a specific job
func (r *jobLicenseRepository) DeleteByJobID(ctx context.Context, jobID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ?", jobID).Delete(&models.JobLicense{}).Error
}

// DeleteByJobAndLicense deletes a specific job license relationship
func (r *jobLicenseRepository) DeleteByJobAndLicense(ctx context.Context, jobID, licenseID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ? AND license_id = ?", jobID, licenseID).Delete(&models.JobLicense{}).Error
}
