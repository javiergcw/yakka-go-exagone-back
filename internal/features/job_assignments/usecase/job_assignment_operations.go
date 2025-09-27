package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_assignments/entity/database"
	"github.com/yakka-backend/internal/features/job_assignments/models"
	"github.com/yakka-backend/internal/features/job_assignments/payload"
)

// JobAssignmentUsecase defines the interface for job assignment business logic
type JobAssignmentUsecase interface {
	CreateAssignment(ctx context.Context, req payload.CreateJobAssignmentRequest) (*models.JobAssignment, error)
	GetAssignmentByID(ctx context.Context, id uuid.UUID) (*models.JobAssignment, error)
	UpdateAssignment(ctx context.Context, id uuid.UUID, req payload.UpdateJobAssignmentRequest) (*models.JobAssignment, error)
	DeleteAssignment(ctx context.Context, id uuid.UUID) error
	GetAssignmentsByJob(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error)
	GetAssignmentsByLabourUser(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error)
	GetAssignmentsWithFilters(ctx context.Context, req payload.GetJobAssignmentsRequest) ([]*models.JobAssignment, int64, error)
	UpdateAssignmentStatus(ctx context.Context, id uuid.UUID, status models.AssignmentStatus) (*models.JobAssignment, error)
	CompleteAssignment(ctx context.Context, id uuid.UUID, req payload.CompleteAssignmentRequest) (*models.JobAssignment, error)
	CancelAssignment(ctx context.Context, id uuid.UUID, req payload.CancelAssignmentRequest) (*models.JobAssignment, error)
}

// JobAssignmentUsecaseImpl implements JobAssignmentUsecase
type JobAssignmentUsecaseImpl struct {
	assignmentRepo database.JobAssignmentRepository
}

// NewJobAssignmentUsecase creates a new job assignment usecase
func NewJobAssignmentUsecase(assignmentRepo database.JobAssignmentRepository) JobAssignmentUsecase {
	return &JobAssignmentUsecaseImpl{
		assignmentRepo: assignmentRepo,
	}
}

// CreateAssignment creates a new job assignment
func (u *JobAssignmentUsecaseImpl) CreateAssignment(ctx context.Context, req payload.CreateJobAssignmentRequest) (*models.JobAssignment, error) {
	// Parse IDs
	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job_id format")
	}

	labourUserID, err := uuid.Parse(req.LabourUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid labour_user_id format")
	}

	applicationID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id format")
	}

	// Check if assignment already exists for this application
	exists, err := u.assignmentRepo.CheckAssignmentExists(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing assignment: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("assignment already exists for this application")
	}

	// Create assignment
	assignment := &models.JobAssignment{
		JobID:         jobID,
		LabourUserID:  labourUserID,
		ApplicationID: applicationID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Status:        models.AssignmentStatusActive,
	}

	if err := u.assignmentRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return assignment, nil
}

// GetAssignmentByID retrieves a job assignment by ID
func (u *JobAssignmentUsecaseImpl) GetAssignmentByID(ctx context.Context, id uuid.UUID) (*models.JobAssignment, error) {
	assignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}
	return assignment, nil
}

// UpdateAssignment updates an existing job assignment
func (u *JobAssignmentUsecaseImpl) UpdateAssignment(ctx context.Context, id uuid.UUID, req payload.UpdateJobAssignmentRequest) (*models.JobAssignment, error) {
	// Get existing assignment
	assignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}

	// Update fields if provided
	if req.StartDate != nil {
		assignment.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		assignment.EndDate = req.EndDate
	}
	if req.Status != nil {
		status := models.AssignmentStatus(*req.Status)
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid status")
		}
		assignment.Status = status
	}

	// Save changes
	if err := u.assignmentRepo.Update(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to update assignment: %w", err)
	}

	return assignment, nil
}

// DeleteAssignment deletes a job assignment
func (u *JobAssignmentUsecaseImpl) DeleteAssignment(ctx context.Context, id uuid.UUID) error {
	// Check if assignment exists
	_, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	if err := u.assignmentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	return nil
}

// GetAssignmentsByJob retrieves all assignments for a specific job
func (u *JobAssignmentUsecaseImpl) GetAssignmentsByJob(ctx context.Context, jobID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	assignments, total, err := u.assignmentRepo.GetByJobID(ctx, jobID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get assignments: %w", err)
	}

	return assignments, total, nil
}

// GetAssignmentsByLabourUser retrieves all assignments for a specific labour user
func (u *JobAssignmentUsecaseImpl) GetAssignmentsByLabourUser(ctx context.Context, labourUserID uuid.UUID, page, limit int) ([]*models.JobAssignment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	assignments, total, err := u.assignmentRepo.GetByLabourUserID(ctx, labourUserID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get assignments: %w", err)
	}

	return assignments, total, nil
}

// GetAssignmentsWithFilters retrieves assignments with multiple filters
func (u *JobAssignmentUsecaseImpl) GetAssignmentsWithFilters(ctx context.Context, req payload.GetJobAssignmentsRequest) ([]*models.JobAssignment, int64, error) {
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
	var jobID, labourUserID, applicationID *uuid.UUID
	var status *models.AssignmentStatus

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

	if req.ApplicationID != nil {
		parsedApplicationID, err := uuid.Parse(*req.ApplicationID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid application_id format")
		}
		applicationID = &parsedApplicationID
	}

	if req.Status != nil {
		parsedStatus := models.AssignmentStatus(*req.Status)
		if !parsedStatus.IsValid() {
			return nil, 0, fmt.Errorf("invalid status")
		}
		status = &parsedStatus
	}

	assignments, total, err := u.assignmentRepo.GetWithFilters(ctx, jobID, labourUserID, applicationID, status, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get assignments: %w", err)
	}

	return assignments, total, nil
}

// UpdateAssignmentStatus updates the status of a job assignment
func (u *JobAssignmentUsecaseImpl) UpdateAssignmentStatus(ctx context.Context, id uuid.UUID, status models.AssignmentStatus) (*models.JobAssignment, error) {
	// Validate status
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid status")
	}

	// Update status
	if err := u.assignmentRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Get updated assignment
	assignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated assignment: %w", err)
	}

	return assignment, nil
}

// CompleteAssignment completes an assignment
func (u *JobAssignmentUsecaseImpl) CompleteAssignment(ctx context.Context, id uuid.UUID, req payload.CompleteAssignmentRequest) (*models.JobAssignment, error) {
	// Check if assignment exists
	assignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}

	// Check if already completed
	if assignment.Status == models.AssignmentStatusCompleted {
		return nil, fmt.Errorf("assignment already completed")
	}

	// Complete assignment
	if err := u.assignmentRepo.CompleteAssignment(ctx, id, req.EndDate); err != nil {
		return nil, fmt.Errorf("failed to complete assignment: %w", err)
	}

	// Get updated assignment
	updatedAssignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated assignment: %w", err)
	}

	return updatedAssignment, nil
}

// CancelAssignment cancels an assignment
func (u *JobAssignmentUsecaseImpl) CancelAssignment(ctx context.Context, id uuid.UUID, req payload.CancelAssignmentRequest) (*models.JobAssignment, error) {
	// Check if assignment exists
	assignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}

	// Check if already cancelled
	if assignment.Status == models.AssignmentStatusCancelled {
		return nil, fmt.Errorf("assignment already cancelled")
	}

	// Cancel assignment
	if err := u.assignmentRepo.CancelAssignment(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to cancel assignment: %w", err)
	}

	// Get updated assignment
	updatedAssignment, err := u.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated assignment: %w", err)
	}

	return updatedAssignment, nil
}

// Helper function to calculate total pages
func calculateTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
