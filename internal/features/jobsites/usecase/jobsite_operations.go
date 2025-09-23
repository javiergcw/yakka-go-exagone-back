package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobsites/entity/database"
	"github.com/yakka-backend/internal/features/jobsites/models"
	"github.com/yakka-backend/internal/features/jobsites/payload"
	"gorm.io/gorm"
)

// JobsiteUsecase defines the interface for jobsite business logic
type JobsiteUsecase interface {
	CreateJobsite(ctx context.Context, req *payload.CreateJobsiteRequest) (*payload.CreateJobsiteResponse, error)
	GetJobsiteByID(ctx context.Context, id uuid.UUID) (*payload.JobsiteResponse, error)
	GetJobsitesByBuilderID(ctx context.Context, builderID uuid.UUID) (*payload.JobsiteListResponse, error)
	UpdateJobsite(ctx context.Context, id uuid.UUID, req *payload.UpdateJobsiteRequest) (*payload.UpdateJobsiteResponse, error)
	DeleteJobsite(ctx context.Context, id uuid.UUID) (*payload.DeleteJobsiteResponse, error)
}

// JobsiteUsecaseImpl implements the JobsiteUsecase interface
type JobsiteUsecaseImpl struct {
	jobsiteRepo database.JobsiteRepository
}

// NewJobsiteUsecase creates a new instance of JobsiteUsecaseImpl
func NewJobsiteUsecaseImpl(jobsiteRepo database.JobsiteRepository) JobsiteUsecase {
	return &JobsiteUsecaseImpl{
		jobsiteRepo: jobsiteRepo,
	}
}

// CreateJobsite creates a new jobsite
func (u *JobsiteUsecaseImpl) CreateJobsite(ctx context.Context, req *payload.CreateJobsiteRequest) (*payload.CreateJobsiteResponse, error) {
	builderID, err := uuid.Parse(req.BuilderID)
	if err != nil {
		return nil, fmt.Errorf("invalid builder ID format: %w", err)
	}

	jobsite := &models.Jobsite{
		BuilderID:   builderID,
		Address:     req.Address,
		City:        req.City,
		Suburb:      req.Suburb,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Phone:       req.Phone,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = u.jobsiteRepo.Create(ctx, jobsite)
	if err != nil {
		return nil, fmt.Errorf("failed to create jobsite: %w", err)
	}

	response := &payload.CreateJobsiteResponse{
		Jobsite: *u.mapToJobsiteResponse(jobsite),
		Message: "Jobsite created successfully",
	}

	return response, nil
}

// GetJobsiteByID retrieves a jobsite by its ID
func (u *JobsiteUsecaseImpl) GetJobsiteByID(ctx context.Context, id uuid.UUID) (*payload.JobsiteResponse, error) {
	jobsite, err := u.jobsiteRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("jobsite not found")
		}
		return nil, fmt.Errorf("failed to retrieve jobsite: %w", err)
	}

	response := u.mapToJobsiteResponse(jobsite)
	return response, nil
}

// GetJobsitesByBuilderID retrieves all jobsites for a specific builder
func (u *JobsiteUsecaseImpl) GetJobsitesByBuilderID(ctx context.Context, builderID uuid.UUID) (*payload.JobsiteListResponse, error) {
	jobsites, err := u.jobsiteRepo.GetByBuilderID(ctx, builderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve jobsites: %w", err)
	}

	jobsiteResponses := make([]payload.JobsiteResponse, len(jobsites))
	for i, jobsite := range jobsites {
		jobsiteResponses[i] = *u.mapToJobsiteResponse(jobsite)
	}

	response := &payload.JobsiteListResponse{
		Jobsites: jobsiteResponses,
		Total:    len(jobsites),
	}

	return response, nil
}

// UpdateJobsite updates an existing jobsite
func (u *JobsiteUsecaseImpl) UpdateJobsite(ctx context.Context, id uuid.UUID, req *payload.UpdateJobsiteRequest) (*payload.UpdateJobsiteResponse, error) {
	jobsite, err := u.jobsiteRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("jobsite not found")
		}
		return nil, fmt.Errorf("failed to retrieve jobsite: %w", err)
	}

	// Update fields if provided
	if req.Address != nil {
		jobsite.Address = *req.Address
	}
	if req.City != nil {
		jobsite.City = *req.City
	}
	if req.Suburb != nil {
		jobsite.Suburb = req.Suburb
	}
	if req.Description != nil {
		jobsite.Description = req.Description
	}
	if req.Latitude != nil {
		jobsite.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		jobsite.Longitude = *req.Longitude
	}
	if req.Phone != nil {
		jobsite.Phone = req.Phone
	}

	jobsite.UpdatedAt = time.Now()

	err = u.jobsiteRepo.Update(ctx, jobsite)
	if err != nil {
		return nil, fmt.Errorf("failed to update jobsite: %w", err)
	}

	response := &payload.UpdateJobsiteResponse{
		Jobsite: *u.mapToJobsiteResponse(jobsite),
		Message: "Jobsite updated successfully",
	}

	return response, nil
}

// DeleteJobsite deletes a jobsite
func (u *JobsiteUsecaseImpl) DeleteJobsite(ctx context.Context, id uuid.UUID) (*payload.DeleteJobsiteResponse, error) {
	// Check if jobsite exists
	_, err := u.jobsiteRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("jobsite not found")
		}
		return nil, fmt.Errorf("failed to retrieve jobsite: %w", err)
	}

	err = u.jobsiteRepo.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete jobsite: %w", err)
	}

	response := &payload.DeleteJobsiteResponse{
		Message: "Jobsite deleted successfully",
	}

	return response, nil
}

// mapToJobsiteResponse maps a Jobsite model to JobsiteResponse
func (u *JobsiteUsecaseImpl) mapToJobsiteResponse(jobsite *models.Jobsite) *payload.JobsiteResponse {
	return &payload.JobsiteResponse{
		ID:          jobsite.ID,
		BuilderID:   jobsite.BuilderID,
		Address:     jobsite.Address,
		City:        jobsite.City,
		Suburb:      jobsite.Suburb,
		Description: jobsite.Description,
		Latitude:    jobsite.Latitude,
		Longitude:   jobsite.Longitude,
		Phone:       jobsite.Phone,
		CreatedAt:   jobsite.CreatedAt,
		UpdatedAt:   jobsite.UpdatedAt,
	}
}
