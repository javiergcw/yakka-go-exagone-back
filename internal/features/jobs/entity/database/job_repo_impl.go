package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
	"gorm.io/gorm"
)

// jobRepository implements JobRepository
type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{
		db: db,
	}
}

// Create creates a new job
func (r *jobRepository) Create(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

// GetByID retrieves a job by ID
func (r *jobRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetByBuilderProfileID retrieves jobs by builder profile ID
func (r *jobRepository) GetByBuilderProfileID(ctx context.Context, builderProfileID uuid.UUID) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.WithContext(ctx).Where("builder_profile_id = ?", builderProfileID).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetByJobsiteID retrieves jobs by jobsite ID
func (r *jobRepository) GetByJobsiteID(ctx context.Context, jobsiteID uuid.UUID) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.WithContext(ctx).Where("jobsite_id = ?", jobsiteID).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetByVisibility retrieves jobs by visibility
func (r *jobRepository) GetByVisibility(ctx context.Context, visibility models.JobVisibility) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.WithContext(ctx).Where("visibility = ?", visibility).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetByVisibilityWithRelations retrieves jobs by visibility with all relations
func (r *jobRepository) GetByVisibilityWithRelations(ctx context.Context, visibility models.JobVisibility) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.WithContext(ctx).
		Preload("JobLicenses").
		Preload("JobSkills").
		Preload("JobRequirements").
		Where("visibility = ?", visibility).
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetAll retrieves all jobs
func (r *jobRepository) GetAll(ctx context.Context) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.WithContext(ctx).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// Update updates a job
func (r *jobRepository) Update(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

// Delete deletes a job by ID
func (r *jobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Job{}).Error
}

// GetWithRelations retrieves a job with all its relations
func (r *jobRepository) GetWithRelations(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := r.db.WithContext(ctx).
		Preload("JobLicenses").
		Preload("JobSkills").
		Preload("JobRequirements").
		Where("id = ?", id).
		First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}
