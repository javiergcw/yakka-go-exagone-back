package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	"github.com/yakka-backend/internal/features/jobs/entity/database"
	"github.com/yakka-backend/internal/features/jobs/models"
	"github.com/yakka-backend/internal/features/jobs/payload"
	jobsite_db "github.com/yakka-backend/internal/features/jobsites/entity/database"
	job_type_db "github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	license_db "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	skill_category_db "github.com/yakka-backend/internal/features/masters/skills/entity/database"
)

// JobUsecase defines the interface for job business logic
type JobUsecase interface {
	CreateJob(ctx context.Context, req payload.CreateJobRequest) (*models.Job, error)
	GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error)
	GetJobsByBuilder(ctx context.Context, builderProfileID uuid.UUID, visibility *models.JobVisibility) ([]*models.Job, error)
	GetJobsByJobsite(ctx context.Context, jobsiteID uuid.UUID, visibility *models.JobVisibility) ([]*models.Job, error)
	GetJobsByVisibility(ctx context.Context, visibility models.JobVisibility) ([]*models.Job, error)
	GetAllJobs(ctx context.Context) ([]*models.Job, error)
	UpdateJob(ctx context.Context, id uuid.UUID, req payload.UpdateJobRequest) (*models.Job, error)
	DeleteJob(ctx context.Context, id uuid.UUID) error
	GetJobWithRelations(ctx context.Context, id uuid.UUID) (*models.Job, error)
}

// jobUsecase implements JobUsecase
type jobUsecase struct {
	jobRepo        database.JobRepository
	jobLicenseRepo database.JobLicenseRepository
	jobSkillRepo   database.JobSkillRepository
	validator      *JobValidationService
}

// NewJobUsecase creates a new job usecase
func NewJobUsecase(
	jobRepo database.JobRepository,
	jobLicenseRepo database.JobLicenseRepository,
	jobSkillRepo database.JobSkillRepository,
	builderRepo builder_db.BuilderProfileRepository,
	jobsiteRepo jobsite_db.JobsiteRepository,
	jobTypeRepo job_type_db.JobTypeRepository,
	licenseRepo license_db.LicenseRepository,
	skillCategoryRepo skill_category_db.SkillCategoryRepository,
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository,
) JobUsecase {
	return &jobUsecase{
		jobRepo:        jobRepo,
		jobLicenseRepo: jobLicenseRepo,
		jobSkillRepo:   jobSkillRepo,
		validator:      NewJobValidationService(builderRepo, jobsiteRepo, jobTypeRepo, licenseRepo, skillCategoryRepo, skillSubcategoryRepo),
	}
}

// CreateJob creates a new job
func (u *jobUsecase) CreateJob(ctx context.Context, req payload.CreateJobRequest) (*models.Job, error) {
	// Validate the request
	if err := u.validator.ValidateCreateJobRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create the job
	job := &models.Job{
		BuilderProfileID:            req.BuilderProfileID,
		JobsiteID:                   req.JobsiteID,
		JobTypeID:                   req.JobTypeID,
		ManyLabours:                 req.ManyLabours,
		OngoingWork:                 req.OngoingWork,
		WageSiteAllowance:           req.WageSiteAllowance,
		WageLeadingHandAllowance:    req.WageLeadingHandAllowance,
		WageProductivityAllowance:   req.WageProductivityAllowance,
		ExtrasOvertimeRate:          req.ExtrasOvertimeRate,
		StartDateWork:               req.StartDateWork,
		EndDateWork:                 req.EndDateWork,
		WorkSaturday:                req.WorkSaturday,
		WorkSunday:                  req.WorkSunday,
		StartTime:                   req.StartTime,
		EndTime:                     req.EndTime,
		Description:                 req.Description,
		PaymentDay:                  req.PaymentDay,
		RequiresSupervisorSignature: req.RequiresSupervisorSignature,
		SupervisorName:              req.SupervisorName,
		Visibility:                  req.Visibility,
		PaymentType:                 req.PaymentType,
	}

	// Create the job
	if err := u.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Create job license relationships
	for _, licenseID := range req.LicenseIDs {
		jobLicense := &models.JobLicense{
			JobID:     job.ID,
			LicenseID: licenseID,
		}
		if err := u.jobLicenseRepo.Create(ctx, jobLicense); err != nil {
			return nil, fmt.Errorf("failed to create job license relationship: %w", err)
		}
	}

	// Create job skill category relationships
	for _, skillCategoryID := range req.SkillCategoryIDs {
		jobSkill := &models.JobSkill{
			JobID:           job.ID,
			SkillCategoryID: &skillCategoryID,
		}
		if err := u.jobSkillRepo.Create(ctx, jobSkill); err != nil {
			return nil, fmt.Errorf("failed to create job skill category relationship: %w", err)
		}
	}

	// Create job skill subcategory relationships
	for _, skillSubcategoryID := range req.SkillSubcategoryIDs {
		jobSkill := &models.JobSkill{
			JobID:              job.ID,
			SkillSubcategoryID: &skillSubcategoryID,
		}
		if err := u.jobSkillRepo.Create(ctx, jobSkill); err != nil {
			return nil, fmt.Errorf("failed to create job skill subcategory relationship: %w", err)
		}
	}

	return job, nil
}

// GetJobByID retrieves a job by ID
func (u *jobUsecase) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}

// GetJobsByBuilder retrieves jobs by builder profile ID
func (u *jobUsecase) GetJobsByBuilder(ctx context.Context, builderProfileID uuid.UUID, visibility *models.JobVisibility) ([]*models.Job, error) {
	jobs, err := u.jobRepo.GetByBuilderProfileID(ctx, builderProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by builder: %w", err)
	}

	// Filter by visibility if provided
	if visibility != nil {
		var filteredJobs []*models.Job
		for _, job := range jobs {
			if job.Visibility == *visibility {
				filteredJobs = append(filteredJobs, job)
			}
		}
		return filteredJobs, nil
	}

	return jobs, nil
}

// GetJobsByJobsite retrieves jobs by jobsite ID
func (u *jobUsecase) GetJobsByJobsite(ctx context.Context, jobsiteID uuid.UUID, visibility *models.JobVisibility) ([]*models.Job, error) {
	jobs, err := u.jobRepo.GetByJobsiteID(ctx, jobsiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by jobsite: %w", err)
	}

	// Filter by visibility if provided
	if visibility != nil {
		var filteredJobs []*models.Job
		for _, job := range jobs {
			if job.Visibility == *visibility {
				filteredJobs = append(filteredJobs, job)
			}
		}
		return filteredJobs, nil
	}

	return jobs, nil
}

// GetJobsByVisibility retrieves jobs by visibility
func (u *jobUsecase) GetJobsByVisibility(ctx context.Context, visibility models.JobVisibility) ([]*models.Job, error) {
	jobs, err := u.jobRepo.GetByVisibility(ctx, visibility)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by visibility: %w", err)
	}
	return jobs, nil
}

// GetAllJobs retrieves all jobs
func (u *jobUsecase) GetAllJobs(ctx context.Context) ([]*models.Job, error) {
	jobs, err := u.jobRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all jobs: %w", err)
	}
	return jobs, nil
}

// UpdateJob updates a job
func (u *jobUsecase) UpdateJob(ctx context.Context, id uuid.UUID, req payload.UpdateJobRequest) (*models.Job, error) {
	// Get existing job
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Update fields if provided
	if req.ManyLabours != nil {
		job.ManyLabours = *req.ManyLabours
	}
	if req.OngoingWork != nil {
		job.OngoingWork = *req.OngoingWork
	}
	if req.WageSiteAllowance != nil {
		job.WageSiteAllowance = req.WageSiteAllowance
	}
	if req.WageLeadingHandAllowance != nil {
		job.WageLeadingHandAllowance = req.WageLeadingHandAllowance
	}
	if req.WageProductivityAllowance != nil {
		job.WageProductivityAllowance = req.WageProductivityAllowance
	}
	if req.ExtrasOvertimeRate != nil {
		job.ExtrasOvertimeRate = req.ExtrasOvertimeRate
	}
	if req.StartDateWork != nil {
		job.StartDateWork = req.StartDateWork
	}
	if req.EndDateWork != nil {
		job.EndDateWork = req.EndDateWork
	}
	if req.WorkSaturday != nil {
		job.WorkSaturday = *req.WorkSaturday
	}
	if req.WorkSunday != nil {
		job.WorkSunday = *req.WorkSunday
	}
	if req.StartTime != nil {
		job.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		job.EndTime = req.EndTime
	}
	if req.Description != nil {
		job.Description = req.Description
	}
	if req.PaymentDay != nil {
		job.PaymentDay = req.PaymentDay
	}
	if req.RequiresSupervisorSignature != nil {
		job.RequiresSupervisorSignature = *req.RequiresSupervisorSignature
	}
	if req.SupervisorName != nil {
		job.SupervisorName = req.SupervisorName
	}
	if req.Visibility != nil {
		job.Visibility = *req.Visibility
	}
	if req.PaymentType != nil {
		job.PaymentType = *req.PaymentType
	}

	// Update the job
	if err := u.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to update job: %w", err)
	}

	// TODO: Handle job licenses and skills relationships updates

	return job, nil
}

// DeleteJob deletes a job
func (u *jobUsecase) DeleteJob(ctx context.Context, id uuid.UUID) error {
	if err := u.jobRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}

// GetJobWithRelations retrieves a job with all its relations
func (u *jobUsecase) GetJobWithRelations(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	job, err := u.jobRepo.GetWithRelations(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job with relations: %w", err)
	}
	return job, nil
}
