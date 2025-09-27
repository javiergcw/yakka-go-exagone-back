package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	job_application_db "github.com/yakka-backend/internal/features/job_applications/entity/database"
	job_application_models "github.com/yakka-backend/internal/features/job_applications/models"
	job_assignment_db "github.com/yakka-backend/internal/features/job_assignments/entity/database"
	job_assignment_models "github.com/yakka-backend/internal/features/job_assignments/models"
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
	GetBuilderApplicants(ctx context.Context, builderProfileID uuid.UUID) ([]payload.JobWithApplicants, error)
	ProcessApplicantDecision(ctx context.Context, builderProfileID uuid.UUID, req payload.BuilderApplicantDecisionRequest) (*payload.BuilderApplicantDecisionResponse, error)
	GetLabourJobs(ctx context.Context, labourUserID uuid.UUID) ([]payload.LabourJobInfo, error)
	ApplyToJob(ctx context.Context, labourUserID uuid.UUID, req payload.LabourApplicationRequest) (*payload.LabourApplicationResponse, error)
}

// jobUsecase implements JobUsecase
type jobUsecase struct {
	jobRepo            database.JobRepository
	jobLicenseRepo     database.JobLicenseRepository
	jobSkillRepo       database.JobSkillRepository
	builderRepo        builder_db.BuilderProfileRepository
	jobsiteRepo        jobsite_db.JobsiteRepository
	jobTypeRepo        job_type_db.JobTypeRepository
	jobApplicationRepo job_application_db.JobApplicationRepository
	jobAssignmentRepo  job_assignment_db.JobAssignmentRepository
	validator          *JobValidationService
}

// NewJobUsecase creates a new job usecase
func NewJobUsecase(
	jobRepo database.JobRepository,
	jobLicenseRepo database.JobLicenseRepository,
	jobSkillRepo database.JobSkillRepository,
	builderRepo builder_db.BuilderProfileRepository,
	jobsiteRepo jobsite_db.JobsiteRepository,
	jobTypeRepo job_type_db.JobTypeRepository,
	jobApplicationRepo job_application_db.JobApplicationRepository,
	jobAssignmentRepo job_assignment_db.JobAssignmentRepository,
	licenseRepo license_db.LicenseRepository,
	skillCategoryRepo skill_category_db.SkillCategoryRepository,
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository,
) JobUsecase {
	return &jobUsecase{
		jobRepo:            jobRepo,
		jobLicenseRepo:     jobLicenseRepo,
		jobSkillRepo:       jobSkillRepo,
		builderRepo:        builderRepo,
		jobsiteRepo:        jobsiteRepo,
		jobTypeRepo:        jobTypeRepo,
		jobApplicationRepo: jobApplicationRepo,
		jobAssignmentRepo:  jobAssignmentRepo,
		validator:          NewJobValidationService(builderRepo, jobsiteRepo, jobTypeRepo, licenseRepo, skillCategoryRepo, skillSubcategoryRepo),
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

// GetBuilderApplicants retrieves all applicants for builder's jobs
func (u *jobUsecase) GetBuilderApplicants(ctx context.Context, builderProfileID uuid.UUID) ([]payload.JobWithApplicants, error) {
	// Get all jobs for this builder
	jobs, err := u.jobRepo.GetByBuilderProfileID(ctx, builderProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builder jobs: %w", err)
	}

	var jobsWithApplicants []payload.JobWithApplicants

	for _, job := range jobs {
		// Get job type for title
		jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID)
		if err != nil {
			continue // Skip jobs with invalid job types
		}

		// Get all applications for this job
		applications, total, err := u.jobApplicationRepo.GetByJobID(ctx, job.ID, 1, 100) // Get up to 100 applications
		if err != nil {
			fmt.Printf("Error getting applications for job %s: %v\n", job.ID, err)
			continue // Skip jobs with errors getting applications
		}

		fmt.Printf("Job %s has %d applications (total: %d)\n", job.ID, len(applications), total)

		// Convert applications to JobApplicantInfo
		applicants := make([]payload.JobApplicantInfo, 0) // Initialize as empty slice, not nil
		for _, app := range applications {
			// Get labour profile information
			// TODO: Get labour profile info from labour_profiles table
			// For now, we'll create a placeholder
			labourInfo := payload.LabourApplicantInfo{
				UserID:    app.LabourUserID.String(),
				FullName:  "Labour User", // TODO: Get from labour profile
				AvatarURL: nil,
				Phone:     nil,
				Email:     "labour@example.com", // TODO: Get from user table
			}

			applicant := payload.JobApplicantInfo{
				ApplicationID: app.ID.String(),
				Status:        string(app.Status),
				CoverLetter:   app.CoverLetter,
				ExpectedRate:  app.ExpectedRate,
				ResumeURL:     app.ResumeURL,
				AppliedAt:     app.CreatedAt,
				Labour:        labourInfo,
			}

			applicants = append(applicants, applicant)
		}

		jobWithApplicants := payload.JobWithApplicants{
			JobID:      job.ID.String(),
			JobTitle:   jobType.Name,
			JobStatus:  "ACTIVE", // TODO: Add status field to job
			CreatedAt:  job.CreatedAt,
			Applicants: applicants,
		}

		jobsWithApplicants = append(jobsWithApplicants, jobWithApplicants)
	}

	return jobsWithApplicants, nil
}

// ProcessApplicantDecision processes hiring or rejection of an applicant
func (u *jobUsecase) ProcessApplicantDecision(ctx context.Context, builderProfileID uuid.UUID, req payload.BuilderApplicantDecisionRequest) (*payload.BuilderApplicantDecisionResponse, error) {
	// Parse application ID
	applicationID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application ID: %w", err)
	}

	// Get the application
	application, err := u.jobApplicationRepo.GetByID(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Get the job to verify it belongs to this builder
	job, err := u.jobRepo.GetByID(ctx, application.JobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Verify the job belongs to this builder
	if job.BuilderProfileID != builderProfileID {
		return nil, fmt.Errorf("application does not belong to this builder")
	}

	response := &payload.BuilderApplicantDecisionResponse{
		ApplicationID: req.ApplicationID,
		Hired:         req.Hired,
		Message:       "Decision processed successfully",
	}

	if req.Hired {
		// Update application status to ACCEPTED
		if err := u.jobApplicationRepo.UpdateStatus(ctx, applicationID, job_application_models.ApplicationStatusAccepted); err != nil {
			return nil, fmt.Errorf("failed to update application status: %w", err)
		}

		// Create job assignment
		assignment := &job_assignment_models.JobAssignment{
			JobID:         application.JobID,
			LabourUserID:  application.LabourUserID,
			ApplicationID: applicationID,
			StartDate:     req.StartDate,
			EndDate:       req.EndDate,
			Status:        job_assignment_models.AssignmentStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Save assignment to database
		if err := u.jobAssignmentRepo.Create(ctx, assignment); err != nil {
			return nil, fmt.Errorf("failed to create job assignment: %w", err)
		}

		assignmentID := assignment.ID.String()
		response.AssignmentID = &assignmentID
	} else {
		// Update application status to REJECTED
		if err := u.jobApplicationRepo.UpdateStatus(ctx, applicationID, job_application_models.ApplicationStatusRejected); err != nil {
			return nil, fmt.Errorf("failed to update application status: %w", err)
		}
	}

	return response, nil
}

// GetLabourJobs retrieves all jobs with application status for a labour user
func (u *jobUsecase) GetLabourJobs(ctx context.Context, labourUserID uuid.UUID) ([]payload.LabourJobInfo, error) {
	// Get all public jobs
	jobs, err := u.jobRepo.GetByVisibility(ctx, models.JobVisibilityPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to get public jobs: %w", err)
	}

	// Convert to LabourJobInfo with application status
	var labourJobs []payload.LabourJobInfo
	for _, job := range jobs {
		// Get builder profile for this job
		builderProfile, err := u.builderRepo.GetByID(ctx, job.BuilderProfileID)
		if err != nil {
			// Skip jobs with invalid builder profiles
			continue
		}

		// Get jobsite information for location
		jobsite, err := u.jobsiteRepo.GetByID(ctx, job.JobsiteID)
		if err != nil {
			// Skip jobs with invalid jobsites
			continue
		}

		// Get job type information
		jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID)
		if err != nil {
			// Skip jobs with invalid job types
			continue
		}

		// Check if labour has applied to this job
		hasApplied := false
		var applicationStatus *string
		var applicationID *string

		// Check if labour has already applied to this job
		application, err := u.jobApplicationRepo.GetByJobAndLabourUser(ctx, job.ID, labourUserID)
		if err == nil && application != nil {
			// User has applied to this job
			hasApplied = true
			status := string(application.Status)
			applicationStatus = &status
			id := application.ID.String()
			applicationID = &id
		}

		labourJob := payload.LabourJobInfo{
			JobID:           job.ID.String(),
			Title:           jobType.Name, // Get from job type
			Description:     getStringValue(job.Description),
			Location:        jobsite.Address, // Get from jobsite
			JobType:         jobType.Name,
			ExperienceLevel: "INTERMEDIATE", // TODO: Get from job requirements
			Status:          "ACTIVE",       // TODO: Add status field to job
			Visibility:      string(job.Visibility),
			Budget:          job.WageSiteAllowance,
			StartDate:       job.StartDateWork,
			EndDate:         job.EndDateWork,
			CreatedAt:       job.CreatedAt,
			UpdatedAt:       job.UpdatedAt,
			Builder: payload.BuilderInfo{
				BuilderID:   builderProfile.ID.String(),
				CompanyName: getStringValue(builderProfile.CompanyName),
				DisplayName: getStringValue(builderProfile.DisplayName),
				Location:    getStringValue(builderProfile.Location),
				AvatarURL:   nil, // TODO: Get from user table
			},
			HasApplied:        hasApplied,
			ApplicationStatus: applicationStatus,
			ApplicationID:     applicationID,
		}

		labourJobs = append(labourJobs, labourJob)
	}

	return labourJobs, nil
}

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ApplyToJob allows a labour user to apply for a job
func (u *jobUsecase) ApplyToJob(ctx context.Context, labourUserID uuid.UUID, req payload.LabourApplicationRequest) (*payload.LabourApplicationResponse, error) {
	// Parse job ID
	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID: %w", err)
	}

	// Check if labour has already applied to this job
	exists, err := u.jobApplicationRepo.CheckApplicationExists(ctx, jobID, labourUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing application: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user has already applied to this job")
	}

	// Get job information for response
	job, err := u.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Get job type for title
	jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID)
	if err != nil {
		return nil, fmt.Errorf("job type not found: %w", err)
	}

	// Create new application
	application := &job_application_models.JobApplication{
		JobID:        jobID,
		LabourUserID: labourUserID,
		Status:       job_application_models.ApplicationStatusApplied,
		CoverLetter:  req.CoverLetter,
		ResumeURL:    req.ResumeURL,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save application to database
	if err := u.jobApplicationRepo.Create(ctx, application); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	response := &payload.LabourApplicationResponse{
		ApplicationID: application.ID.String(),
		JobID:         req.JobID,
		JobTitle:      jobType.Name,
		Status:        string(application.Status),
		CoverLetter:   application.CoverLetter,
		ResumeURL:     application.ResumeURL,
		AppliedAt:     application.CreatedAt,
		Message:       "Application submitted successfully",
	}

	return response, nil
}
