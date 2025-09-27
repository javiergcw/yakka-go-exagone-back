package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_applications/entity/database"
	"github.com/yakka-backend/internal/features/job_applications/models"
	"github.com/yakka-backend/internal/features/job_applications/payload"
)

// JobApplicationUsecase defines the interface for job application business logic
type JobApplicationUsecase interface {
	CreateApplication(ctx context.Context, labourUserID uuid.UUID, req payload.CreateJobApplicationRequest) (*models.JobApplication, error)
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*models.JobApplication, error)
	UpdateApplication(ctx context.Context, id uuid.UUID, req payload.UpdateJobApplicationRequest) (*models.JobApplication, error)
	DeleteApplication(ctx context.Context, id uuid.UUID) error
	GetApplicationsByJob(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error)
	GetApplicationsByLabourUser(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error)
	GetApplicationsWithFilters(ctx context.Context, req payload.GetJobApplicationsRequest) ([]*models.JobApplication, int64, error)
	UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status models.ApplicationStatus) (*models.JobApplication, error)
	WithdrawApplication(ctx context.Context, id uuid.UUID) error
}

// JobApplicationUsecaseImpl implements JobApplicationUsecase
type JobApplicationUsecaseImpl struct {
	applicationRepo database.JobApplicationRepository
}

// NewJobApplicationUsecase creates a new job application usecase
func NewJobApplicationUsecase(applicationRepo database.JobApplicationRepository) JobApplicationUsecase {
	return &JobApplicationUsecaseImpl{
		applicationRepo: applicationRepo,
	}
}

// CreateApplication creates a new job application
func (u *JobApplicationUsecaseImpl) CreateApplication(ctx context.Context, labourUserID uuid.UUID, req payload.CreateJobApplicationRequest) (*models.JobApplication, error) {
	// Parse job ID
	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job_id format")
	}

	// Check if application already exists
	exists, err := u.applicationRepo.CheckApplicationExists(ctx, jobID, labourUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing application: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("application already exists for this job")
	}

	// Create application
	application := &models.JobApplication{
		JobID:        jobID,
		LabourUserID: labourUserID,
		Status:       models.ApplicationStatusApplied,
		CoverLetter:  req.CoverLetter,
		ExpectedRate: req.ExpectedRate,
		ResumeURL:    req.ResumeURL,
	}

	if err := u.applicationRepo.Create(ctx, application); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return application, nil
}

// GetApplicationByID retrieves a job application by ID
func (u *JobApplicationUsecaseImpl) GetApplicationByID(ctx context.Context, id uuid.UUID) (*models.JobApplication, error) {
	application, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return application, nil
}

// UpdateApplication updates an existing job application
func (u *JobApplicationUsecaseImpl) UpdateApplication(ctx context.Context, id uuid.UUID, req payload.UpdateJobApplicationRequest) (*models.JobApplication, error) {
	// Get existing application
	application, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Update fields if provided
	if req.Status != nil {
		application.Status = *req.Status
	}
	if req.CoverLetter != nil {
		application.CoverLetter = req.CoverLetter
	}
	if req.ExpectedRate != nil {
		application.ExpectedRate = req.ExpectedRate
	}
	if req.ResumeURL != nil {
		application.ResumeURL = req.ResumeURL
	}

	// Save changes
	if err := u.applicationRepo.Update(ctx, application); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return application, nil
}

// DeleteApplication deletes a job application
func (u *JobApplicationUsecaseImpl) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	// Check if application exists
	_, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	if err := u.applicationRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}

// GetApplicationsByJob retrieves all applications for a specific job
func (u *JobApplicationUsecaseImpl) GetApplicationsByJob(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	applications, total, err := u.applicationRepo.GetByJobID(ctx, jobID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
	}

	return applications, total, nil
}

// GetApplicationsByLabourUser retrieves all applications for a specific labour user
func (u *JobApplicationUsecaseImpl) GetApplicationsByLabourUser(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	applications, total, err := u.applicationRepo.GetByLabourUserID(ctx, labourUserID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
	}

	return applications, total, nil
}

// GetApplicationsWithFilters retrieves applications with multiple filters
func (u *JobApplicationUsecaseImpl) GetApplicationsWithFilters(ctx context.Context, req payload.GetJobApplicationsRequest) ([]*models.JobApplication, int64, error) {
	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse filters
	var jobID, labourUserID *uuid.UUID
	var status *models.ApplicationStatus

	if req.JobID != nil {
		parsedJobID, err := uuid.Parse(*req.JobID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid job_id format")
		}
		jobID = &parsedJobID
	}

	if req.LabourUserID != nil {
		parsedLabourUserID, err := uuid.Parse(*req.LabourUserID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid labour_user_id format")
		}
		labourUserID = &parsedLabourUserID
	}

	if req.Status != nil {
		parsedStatus := models.ApplicationStatus(*req.Status)
		if !parsedStatus.IsValid() {
			return nil, 0, fmt.Errorf("invalid status")
		}
		status = &parsedStatus
	}

	applications, total, err := u.applicationRepo.GetWithFilters(ctx, jobID, labourUserID, status, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
	}

	return applications, total, nil
}

// UpdateApplicationStatus updates the status of a job application
func (u *JobApplicationUsecaseImpl) UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status models.ApplicationStatus) (*models.JobApplication, error) {
	// Validate status
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid status")
	}

	// Update status
	if err := u.applicationRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Get updated application
	application, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated application: %w", err)
	}

	return application, nil
}

// WithdrawApplication withdraws an application
func (u *JobApplicationUsecaseImpl) WithdrawApplication(ctx context.Context, id uuid.UUID) error {
	// Check if application exists
	application, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if already withdrawn
	if application.Status == models.ApplicationStatusWithdrawn {
		return fmt.Errorf("application already withdrawn")
	}

	// Withdraw application
	if err := u.applicationRepo.WithdrawApplication(ctx, id); err != nil {
		return fmt.Errorf("failed to withdraw application: %w", err)
	}

	return nil
}

// Helper function to calculate total pages
func calculateTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
