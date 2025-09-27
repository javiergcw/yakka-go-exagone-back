package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_applications/models"
	"gorm.io/gorm"
)

// JobApplicationRepositoryImpl implements JobApplicationRepository
type JobApplicationRepositoryImpl struct {
	db *gorm.DB
}

// NewJobApplicationRepository creates a new job application repository
func NewJobApplicationRepository(db *gorm.DB) JobApplicationRepository {
	return &JobApplicationRepositoryImpl{db: db}
}

// Create creates a new job application
func (r *JobApplicationRepositoryImpl) Create(ctx context.Context, application *models.JobApplication) error {
	return r.db.WithContext(ctx).Create(application).Error
}

// GetByID retrieves a job application by ID
func (r *JobApplicationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.JobApplication, error) {
	var application models.JobApplication
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

// GetByJobAndLabourUser retrieves a job application by job ID and labour user ID
func (r *JobApplicationRepositoryImpl) GetByJobAndLabourUser(ctx context.Context, jobID, labourUserID uuid.UUID) (*models.JobApplication, error) {
	var application models.JobApplication
	err := r.db.WithContext(ctx).Where("job_id = ? AND labour_user_id = ?", jobID, labourUserID).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

// Update updates an existing job application
func (r *JobApplicationRepositoryImpl) Update(ctx context.Context, application *models.JobApplication) error {
	application.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(application).Error
}

// Delete deletes a job application
func (r *JobApplicationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.JobApplication{}).Error
}

// GetByJobID retrieves all applications for a specific job
func (r *JobApplicationRepositoryImpl) GetByJobID(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobApplication{}).Where("job_id = ?", jobID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&applications).Error
	return applications, total, err
}

// GetByLabourUserID retrieves all applications for a specific labour user
func (r *JobApplicationRepositoryImpl) GetByLabourUserID(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobApplication{}).Where("labour_user_id = ?", labourUserID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&applications).Error
	return applications, total, err
}

// GetByStatus retrieves applications by status
func (r *JobApplicationRepositoryImpl) GetByStatus(ctx context.Context, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobApplication{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&applications).Error
	return applications, total, err
}

// GetWithFilters retrieves applications with multiple filters
func (r *JobApplicationRepositoryImpl) GetWithFilters(ctx context.Context, jobID, labourUserID *uuid.UUID, status *models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobApplication{})

	// Apply filters
	if jobID != nil {
		query = query.Where("job_id = ?", *jobID)
	}
	if labourUserID != nil {
		query = query.Where("labour_user_id = ?", *labourUserID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&applications).Error
	return applications, total, err
}

// UpdateStatus updates the status of a job application
func (r *JobApplicationRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ApplicationStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// If withdrawing, set withdrawn_at
	if status == models.ApplicationStatusWithdrawn {
		updates["withdrawn_at"] = time.Now()
	}

	return r.db.WithContext(ctx).Model(&models.JobApplication{}).Where("id = ?", id).Updates(updates).Error
}

// WithdrawApplication withdraws an application
func (r *JobApplicationRepositoryImpl) WithdrawApplication(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       models.ApplicationStatusWithdrawn,
		"withdrawn_at": now,
		"updated_at":   now,
	}

	return r.db.WithContext(ctx).Model(&models.JobApplication{}).Where("id = ?", id).Updates(updates).Error
}

// CheckApplicationExists checks if an application already exists for a job and user
func (r *JobApplicationRepositoryImpl) CheckApplicationExists(ctx context.Context, jobID, labourUserID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.JobApplication{}).
		Where("job_id = ? AND labour_user_id = ?", jobID, labourUserID).
		Count(&count).Error

	return count > 0, err
}
