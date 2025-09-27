package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_assignments/models"
	"gorm.io/gorm"
)

// JobAssignmentRepositoryImpl implements JobAssignmentRepository
type JobAssignmentRepositoryImpl struct {
	db *gorm.DB
}

// NewJobAssignmentRepository creates a new job assignment repository
func NewJobAssignmentRepository(db *gorm.DB) JobAssignmentRepository {
	return &JobAssignmentRepositoryImpl{db: db}
}

// Create creates a new job assignment
func (r *JobAssignmentRepositoryImpl) Create(ctx context.Context, assignment *models.JobAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

// GetByID retrieves a job assignment by ID
func (r *JobAssignmentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.JobAssignment, error) {
	var assignment models.JobAssignment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetByApplicationID retrieves a job assignment by application ID
func (r *JobAssignmentRepositoryImpl) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) (*models.JobAssignment, error) {
	var assignment models.JobAssignment
	err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// Update updates an existing job assignment
func (r *JobAssignmentRepositoryImpl) Update(ctx context.Context, assignment *models.JobAssignment) error {
	assignment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(assignment).Error
}

// Delete deletes a job assignment
func (r *JobAssignmentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.JobAssignment{}).Error
}

// GetByJobID retrieves all assignments for a specific job
func (r *JobAssignmentRepositoryImpl) GetByJobID(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error) {
	var assignments []*models.JobAssignment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobAssignment{}).Where("job_id = ?", jobID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&assignments).Error
	return assignments, total, err
}

// GetByLabourUserID retrieves all assignments for a specific labour user
func (r *JobAssignmentRepositoryImpl) GetByLabourUserID(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error) {
	var assignments []*models.JobAssignment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobAssignment{}).Where("labour_user_id = ?", labourUserID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&assignments).Error
	return assignments, total, err
}

// GetByStatus retrieves assignments by status
func (r *JobAssignmentRepositoryImpl) GetByStatus(ctx context.Context, status models.AssignmentStatus, page, limit int) ([]*models.JobAssignment, int64, error) {
	var assignments []*models.JobAssignment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobAssignment{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&assignments).Error
	return assignments, total, err
}

// GetWithFilters retrieves assignments with multiple filters
func (r *JobAssignmentRepositoryImpl) GetWithFilters(ctx context.Context, jobID, labourUserID, applicationID *uuid.UUID, status *models.AssignmentStatus, page, limit int) ([]*models.JobAssignment, int64, error) {
	var assignments []*models.JobAssignment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobAssignment{})

	// Apply filters
	if jobID != nil {
		query = query.Where("job_id = ?", *jobID)
	}
	if labourUserID != nil {
		query = query.Where("labour_user_id = ?", *labourUserID)
	}
	if applicationID != nil {
		query = query.Where("application_id = ?", *applicationID)
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
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&assignments).Error
	return assignments, total, err
}

// UpdateStatus updates the status of a job assignment
func (r *JobAssignmentRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.AssignmentStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	return r.db.WithContext(ctx).Model(&models.JobAssignment{}).Where("id = ?", id).Updates(updates).Error
}

// CompleteAssignment completes an assignment
func (r *JobAssignmentRepositoryImpl) CompleteAssignment(ctx context.Context, id uuid.UUID, endDate *time.Time) error {
	updates := map[string]interface{}{
		"status":     models.AssignmentStatusCompleted,
		"updated_at": time.Now(),
	}

	if endDate != nil {
		updates["end_date"] = *endDate
	} else {
		updates["end_date"] = time.Now()
	}

	return r.db.WithContext(ctx).Model(&models.JobAssignment{}).Where("id = ?", id).Updates(updates).Error
}

// CancelAssignment cancels an assignment
func (r *JobAssignmentRepositoryImpl) CancelAssignment(ctx context.Context, id uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     models.AssignmentStatusCancelled,
		"updated_at": time.Now(),
	}

	return r.db.WithContext(ctx).Model(&models.JobAssignment{}).Where("id = ?", id).Updates(updates).Error
}

// CheckAssignmentExists checks if an assignment already exists for an application
func (r *JobAssignmentRepositoryImpl) CheckAssignmentExists(ctx context.Context, applicationID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.JobAssignment{}).
		Where("application_id = ?", applicationID).
		Count(&count).Error

	return count > 0, err
}
