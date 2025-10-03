package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	auth_user_db "github.com/yakka-backend/internal/features/auth/user/entity/database"
	auth_user_models "github.com/yakka-backend/internal/features/auth/user/models"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	builder_models "github.com/yakka-backend/internal/features/builder_profiles/models"
	job_application_db "github.com/yakka-backend/internal/features/job_applications/entity/database"
	job_application_models "github.com/yakka-backend/internal/features/job_applications/models"
	job_assignment_db "github.com/yakka-backend/internal/features/job_assignments/entity/database"
	job_assignment_models "github.com/yakka-backend/internal/features/job_assignments/models"
	"github.com/yakka-backend/internal/features/jobs/entity/database"
	"github.com/yakka-backend/internal/features/jobs/models"
	"github.com/yakka-backend/internal/features/jobs/payload"
	jobsite_db "github.com/yakka-backend/internal/features/jobsites/entity/database"
	jobsite_models "github.com/yakka-backend/internal/features/jobsites/models"
	job_requirement_db "github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	job_type_db "github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	job_type_models "github.com/yakka-backend/internal/features/masters/job_types/models"
	license_db "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	skill_category_db "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	"gorm.io/gorm"
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
	GetBuilderApplicantsByJobsite(ctx context.Context, builderProfileID uuid.UUID) ([]payload.JobsiteWithJobs, error)
	ProcessApplicantDecision(ctx context.Context, builderProfileID uuid.UUID, req payload.BuilderApplicantDecisionRequest) (*payload.BuilderApplicantDecisionResponse, error)
	GetLabourJobs(ctx context.Context, labourUserID uuid.UUID) ([]payload.LabourJobInfo, error)
	ApplyToJob(ctx context.Context, labourUserID uuid.UUID, req payload.LabourApplicationRequest) (*payload.LabourApplicationResponse, error)
	GetLabourJobDetail(ctx context.Context, jobID uuid.UUID, labourUserID uuid.UUID) (*payload.LabourJobDetailResponse, error)
	GetBuilderJobDetail(ctx context.Context, jobID uuid.UUID, builderProfileID uuid.UUID) (*payload.GetJobResponse, error)
	UpdateJobVisibility(ctx context.Context, jobID uuid.UUID, builderProfileID uuid.UUID, req payload.UpdateJobVisibilityRequest) (*payload.UpdateJobVisibilityResponse, error)
	GetLabourApplicants(ctx context.Context, labourUserID uuid.UUID) (*payload.LabourApplicantsResponse, error)
}

// jobUsecase implements JobUsecase
type jobUsecase struct {
	jobRepo               database.JobRepository
	jobLicenseRepo        database.JobLicenseRepository
	jobSkillRepo          database.JobSkillRepository
	jobJobRequirementRepo database.JobJobRequirementRepository
	jobRequirementRepo    job_requirement_db.JobRequirementRepository
	builderRepo           builder_db.BuilderProfileRepository
	jobsiteRepo           jobsite_db.JobsiteRepository
	jobTypeRepo           job_type_db.JobTypeRepository
	jobApplicationRepo    job_application_db.JobApplicationRepository
	jobAssignmentRepo     job_assignment_db.JobAssignmentRepository
	licenseRepo           license_db.LicenseRepository
	skillCategoryRepo     skill_category_db.SkillCategoryRepository
	skillSubcategoryRepo  skill_category_db.SkillSubcategoryRepository
	userRepo              auth_user_db.UserRepository
	validator             *JobValidationService
}

// NewJobUsecase creates a new job usecase
func NewJobUsecase(
	jobRepo database.JobRepository,
	jobLicenseRepo database.JobLicenseRepository,
	jobSkillRepo database.JobSkillRepository,
	jobJobRequirementRepo database.JobJobRequirementRepository,
	jobRequirementRepo job_requirement_db.JobRequirementRepository,
	builderRepo builder_db.BuilderProfileRepository,
	jobsiteRepo jobsite_db.JobsiteRepository,
	jobTypeRepo job_type_db.JobTypeRepository,
	jobApplicationRepo job_application_db.JobApplicationRepository,
	jobAssignmentRepo job_assignment_db.JobAssignmentRepository,
	licenseRepo license_db.LicenseRepository,
	skillCategoryRepo skill_category_db.SkillCategoryRepository,
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository,
	userRepo auth_user_db.UserRepository,
) JobUsecase {
	return &jobUsecase{
		jobRepo:               jobRepo,
		jobLicenseRepo:        jobLicenseRepo,
		jobSkillRepo:          jobSkillRepo,
		jobJobRequirementRepo: jobJobRequirementRepo,
		jobRequirementRepo:    jobRequirementRepo,
		builderRepo:           builderRepo,
		jobsiteRepo:           jobsiteRepo,
		jobTypeRepo:           jobTypeRepo,
		jobApplicationRepo:    jobApplicationRepo,
		jobAssignmentRepo:     jobAssignmentRepo,
		licenseRepo:           licenseRepo,
		skillCategoryRepo:     skillCategoryRepo,
		skillSubcategoryRepo:  skillSubcategoryRepo,
		userRepo:              userRepo,
		validator:             NewJobValidationService(builderRepo, jobsiteRepo, jobTypeRepo, licenseRepo, skillCategoryRepo, skillSubcategoryRepo, jobRequirementRepo),
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
		WageHourlyRate:              req.WageHourlyRate,
		TravelAllowance:             req.TravelAllowance,
		GST:                         req.GST,
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

	// Create job skill relationships
	// Handle new format: JobSkills with category and subcategory together
	log.Printf("🔍 CreateJob - JobSkills count: %d", len(req.JobSkills))
	if len(req.JobSkills) > 0 {
		for i, jobSkillReq := range req.JobSkills {
			log.Printf("🔍 CreateJob - Creating JobSkill %d: CategoryID=%v, SubcategoryID=%v", i, jobSkillReq.SkillCategoryID, jobSkillReq.SkillSubcategoryID)
			jobSkill := &models.JobSkill{
				JobID:              job.ID,
				SkillCategoryID:    jobSkillReq.SkillCategoryID,
				SkillSubcategoryID: jobSkillReq.SkillSubcategoryID,
			}
			if err := u.jobSkillRepo.Create(ctx, jobSkill); err != nil {
				log.Printf("🚫 CreateJob - Failed to create job skill relationship: %v", err)
				return nil, fmt.Errorf("failed to create job skill relationship: %w", err)
			}
			log.Printf("🔍 CreateJob - JobSkill created successfully: ID=%s", jobSkill.ID)
		}
	} else {
		// Handle legacy format: separate arrays for backward compatibility
		// Create records for skill categories only (without subcategories)
		for _, skillCategoryID := range req.SkillCategoryIDs {
			jobSkill := &models.JobSkill{
				JobID:           job.ID,
				SkillCategoryID: &skillCategoryID,
			}
			if err := u.jobSkillRepo.Create(ctx, jobSkill); err != nil {
				return nil, fmt.Errorf("failed to create job skill category relationship: %w", err)
			}
		}

		// Create records for skill subcategories only (without categories)
		for _, skillSubcategoryID := range req.SkillSubcategoryIDs {
			jobSkill := &models.JobSkill{
				JobID:              job.ID,
				SkillSubcategoryID: &skillSubcategoryID,
			}
			if err := u.jobSkillRepo.Create(ctx, jobSkill); err != nil {
				return nil, fmt.Errorf("failed to create job skill subcategory relationship: %w", err)
			}
		}
	}

	// Create job requirement relationships
	for _, jobRequirementID := range req.JobRequirementIDs {
		jobJobRequirement := &models.JobJobRequirement{
			JobID:            job.ID,
			JobRequirementID: jobRequirementID,
		}
		if err := u.jobJobRequirementRepo.Create(ctx, jobJobRequirement); err != nil {
			return nil, fmt.Errorf("failed to create job requirement relationship: %w", err)
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
	if req.WageHourlyRate != nil {
		job.WageHourlyRate = req.WageHourlyRate
	}
	if req.TravelAllowance != nil {
		job.TravelAllowance = req.TravelAllowance
	}
	if req.GST != nil {
		job.GST = req.GST
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
			// Get labour user information from users table
			labourUser, err := u.userRepo.GetByID(ctx, app.LabourUserID)
			if err != nil {
				log.Printf("Error getting labour user %s: %v", app.LabourUserID, err)
				// Create a fallback labour info if user not found
				labourInfo := payload.LabourApplicantInfo{
					UserID:    app.LabourUserID.String(),
					FullName:  "Usuario Labour",
					AvatarURL: nil,
					Phone:     nil,
					Email:     "usuario@ejemplo.com",
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
				continue
			}

			// Build labour info from user data
			labourInfo := u.buildLabourApplicantInfoFromUser(labourUser)

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

// GetBuilderApplicantsByJobsite retrieves all applicants for builder's jobs grouped by jobsite
func (u *jobUsecase) GetBuilderApplicantsByJobsite(ctx context.Context, builderProfileID uuid.UUID) ([]payload.JobsiteWithJobs, error) {
	// Get all jobs for this builder
	jobs, err := u.jobRepo.GetByBuilderProfileID(ctx, builderProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builder jobs: %w", err)
	}

	// Group jobs by jobsite
	jobsiteMap := make(map[uuid.UUID]*payload.JobsiteWithJobs)

	for _, job := range jobs {
		// Get jobsite information
		jobsite, err := u.jobsiteRepo.GetByID(ctx, job.JobsiteID)
		if err != nil {
			continue // Skip jobs with invalid jobsites
		}

		// Get job type for title
		jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID)
		if err != nil {
			continue // Skip jobs with invalid job types
		}

		// Get all applications for this job
		applications, _, err := u.jobApplicationRepo.GetByJobID(ctx, job.ID, 1, 100) // Get up to 100 applications
		if err != nil {
			fmt.Printf("Error getting applications for job %s: %v\n", job.ID, err)
			continue // Skip jobs with errors getting applications
		}

		// Convert applications to JobApplicantInfo
		applicants := make([]payload.JobApplicantInfo, 0)
		for _, app := range applications {
			// Get labour user information from users table
			labourUser, err := u.userRepo.GetByID(ctx, app.LabourUserID)
			if err != nil {
				log.Printf("Error getting labour user %s: %v", app.LabourUserID, err)
				// Create a fallback labour info if user not found
				labourInfo := payload.LabourApplicantInfo{
					UserID:    app.LabourUserID.String(),
					FullName:  "Usuario Labour",
					AvatarURL: nil,
					Phone:     nil,
					Email:     "usuario@ejemplo.com",
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
				continue
			}

			// Build labour info from user data
			labourInfo := u.buildLabourApplicantInfoFromUser(labourUser)

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

		// Create or get jobsite entry
		if jobsiteEntry, exists := jobsiteMap[jobsite.ID]; exists {
			// Add job to existing jobsite
			jobWithApplicants := payload.JobWithApplicants{
				JobID:      job.ID.String(),
				JobTitle:   jobType.Name,
				JobStatus:  "ACTIVE", // TODO: Add status field to job
				CreatedAt:  job.CreatedAt,
				Applicants: applicants,
			}
			jobsiteEntry.Jobs = append(jobsiteEntry.Jobs, jobWithApplicants)
		} else {
			// Create new jobsite entry
			jobWithApplicants := payload.JobWithApplicants{
				JobID:      job.ID.String(),
				JobTitle:   jobType.Name,
				JobStatus:  "ACTIVE", // TODO: Add status field to job
				CreatedAt:  job.CreatedAt,
				Applicants: applicants,
			}

			jobsiteEntry := &payload.JobsiteWithJobs{
				JobsiteID:   jobsite.ID.String(),
				JobsiteName: getStringValue(jobsite.Description), // Use description as name
				Address:     jobsite.Address,
				City:        jobsite.City,
				Suburb:      jobsite.Suburb,
				Jobs:        []payload.JobWithApplicants{jobWithApplicants},
			}
			jobsiteMap[jobsite.ID] = jobsiteEntry
		}
	}

	// Convert map to slice
	var jobsitesWithJobs []payload.JobsiteWithJobs
	for _, jobsite := range jobsiteMap {
		jobsitesWithJobs = append(jobsitesWithJobs, *jobsite)
	}

	return jobsitesWithJobs, nil
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
		Hired:         *req.Hired,
		Message:       "Decision processed successfully",
	}

	if *req.Hired {
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
	// Get all public jobs with relations
	jobs, err := u.jobRepo.GetByVisibilityWithRelations(ctx, models.JobVisibilityPublic)
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

		// Calculate total wage for budget
		totalWage := u.calculateTotalWage(job)

		// Create jobsite info
		jobsiteInfo := &payload.JobsiteInfo{
			ID:          jobsite.ID.String(),
			Name:        getStringValue(jobsite.Description), // Use description as name
			Address:     jobsite.Address,
			City:        jobsite.City,
			Suburb:      jobsite.Suburb,
			Description: jobsite.Description,
			Latitude:    jobsite.Latitude,
			Longitude:   jobsite.Longitude,
			Phone:       jobsite.Phone,
			CreatedAt:   jobsite.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   jobsite.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Create skills info
		var skillsInfo []payload.JobSkillInfo
		for _, jobSkill := range job.JobSkills {
			skillInfo := payload.JobSkillInfo{
				ID:                 jobSkill.ID.String(),
				JobID:              jobSkill.JobID.String(),
				SkillCategoryID:    nil,
				SkillSubcategoryID: nil,
				CreatedAt:          jobSkill.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}

			// Get skill category details if available
			if jobSkill.SkillCategoryID != nil {
				categoryID := jobSkill.SkillCategoryID.String()
				skillInfo.SkillCategoryID = &categoryID

				skillCategory, err := u.skillCategoryRepo.GetByID(ctx, *jobSkill.SkillCategoryID)
				if err == nil {
					skillInfo.SkillCategory = &payload.SkillCategoryInfo{
						ID:          skillCategory.ID.String(),
						Name:        skillCategory.Name,
						Description: &skillCategory.Description,
					}
				}
			}

			// Get skill subcategory details if available
			if jobSkill.SkillSubcategoryID != nil {
				subcategoryID := jobSkill.SkillSubcategoryID.String()
				skillInfo.SkillSubcategoryID = &subcategoryID

				skillSubcategory, err := u.skillSubcategoryRepo.GetByID(ctx, *jobSkill.SkillSubcategoryID)
				if err == nil {
					skillInfo.SkillSubcategory = &payload.SkillSubcategoryInfo{
						ID:          skillSubcategory.ID.String(),
						Name:        skillSubcategory.Name,
						Description: &skillSubcategory.Description,
					}
				}
			}

			skillsInfo = append(skillsInfo, skillInfo)
		}

		labourJob := payload.LabourJobInfo{
			JobID:           job.ID.String(),
			Title:           jobType.Name, // Get from job type
			Description:     getStringValue(job.Description),
			Location:        jobsite.Address, // Get from jobsite
			JobType:         jobType.Name,
			ManyLabours:     job.ManyLabours,
			ExperienceLevel: "INTERMEDIATE", // TODO: Get from job requirements
			Status:          "ACTIVE",       // TODO: Add status field to job
			Visibility:      string(job.Visibility),
			TotalWage:       &totalWage,
			StartDate:       job.StartDateWork,
			EndDate:         job.EndDateWork,
			CreatedAt:       job.CreatedAt,
			UpdatedAt:       job.UpdatedAt,
			Builder: payload.BuilderInfo{
				BuilderID:   builderProfile.ID.String(),
				CompanyName: getCompanyName(builderProfile.Company),
				DisplayName: getStringValue(builderProfile.DisplayName),
				Location:    getStringValue(builderProfile.Location),
				AvatarURL:   nil, // TODO: Get from user table
			},
			Jobsite:           jobsiteInfo,
			Skills:            skillsInfo,
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

// getCompanyName safely gets company name from Company model
func getCompanyName(company *builder_models.Company) string {
	if company == nil {
		return ""
	}
	return company.Name
}

// buildLabourApplicantInfoFromUser builds LabourApplicantInfo from User model
func (u *jobUsecase) buildLabourApplicantInfoFromUser(user *auth_user_models.User) payload.LabourApplicantInfo {
	// Build full name from first and last name
	var fullName string
	if user.FirstName != nil && user.LastName != nil {
		fullName = *user.FirstName + " " + *user.LastName
	} else if user.FirstName != nil {
		fullName = *user.FirstName
	} else if user.LastName != nil {
		fullName = *user.LastName
	} else {
		fullName = "Usuario Labour" // Fallback
	}

	return payload.LabourApplicantInfo{
		UserID:    user.ID.String(),
		FullName:  fullName,
		AvatarURL: user.Photo,
		Phone:     user.Phone,
		Email:     user.Email,
	}
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

// GetLabourJobDetail retrieves a job detail for a labour user with application info
func (u *jobUsecase) GetLabourJobDetail(ctx context.Context, jobID uuid.UUID, labourUserID uuid.UUID) (*payload.LabourJobDetailResponse, error) {
	// Get job with all relations
	job, err := u.jobRepo.GetWithRelations(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Get additional relation data
	builderProfile, err := u.builderRepo.GetByID(ctx, job.BuilderProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builder profile: %w", err)
	}

	jobsite, err := u.jobsiteRepo.GetByID(ctx, job.JobsiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobsite: %w", err)
	}

	jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job type: %w", err)
	}

	// Convert job to response with additional data
	jobResp := u.convertToJobResponseWithRelations(ctx, job, builderProfile, jobsite, jobType)

	// Check if user has applied to this job
	application, err := u.jobApplicationRepo.GetByJobAndLabourUser(ctx, jobID, labourUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check application: %w", err)
	}

	response := &payload.LabourJobDetailResponse{
		Job:     jobResp,
		Message: "Job detail retrieved successfully",
	}

	// If application exists, add it to response
	if application != nil {
		applicationInfo := &payload.JobApplicationInfo{
			ID:           application.ID,
			Status:       string(application.Status),
			CoverLetter:  application.CoverLetter,
			ExpectedRate: application.ExpectedRate,
			ResumeURL:    application.ResumeURL,
			CreatedAt:    application.CreatedAt,
			UpdatedAt:    application.UpdatedAt,
			WithdrawnAt:  application.WithdrawnAt,
		}
		response.Application = applicationInfo
	}

	return response, nil
}

// convertToJobResponseWithRelations converts a Job model to JobResponse with full relations
func (u *jobUsecase) convertToJobResponseWithRelations(ctx context.Context, job *models.Job, builderProfile *builder_models.BuilderProfile, jobsite *jobsite_models.Jobsite, jobType *job_type_models.JobType) payload.JobResponse {
	jobResp := payload.JobResponse{
		ID:                          job.ID,
		ManyLabours:                 job.ManyLabours,
		OngoingWork:                 job.OngoingWork,
		WageSiteAllowance:           job.WageSiteAllowance,
		WageLeadingHandAllowance:    job.WageLeadingHandAllowance,
		WageProductivityAllowance:   job.WageProductivityAllowance,
		ExtrasOvertimeRate:          job.ExtrasOvertimeRate,
		WageHourlyRate:              job.WageHourlyRate,
		TravelAllowance:             job.TravelAllowance,
		GST:                         job.GST,
		StartDateWork:               job.StartDateWork,
		EndDateWork:                 job.EndDateWork,
		WorkSaturday:                job.WorkSaturday,
		WorkSunday:                  job.WorkSunday,
		StartTime:                   job.StartTime,
		EndTime:                     job.EndTime,
		Description:                 job.Description,
		PaymentDay:                  job.PaymentDay,
		RequiresSupervisorSignature: job.RequiresSupervisorSignature,
		SupervisorName:              job.SupervisorName,
		Visibility:                  job.Visibility,
		PaymentType:                 job.PaymentType,
		CreatedAt:                   job.CreatedAt,
		UpdatedAt:                   job.UpdatedAt,
	}

	// Calculate total wage
	totalWage := u.calculateTotalWage(job)
	jobResp.TotalWage = &totalWage

	// Add builder profile information
	if builderProfile != nil {
		jobResp.BuilderProfile = &payload.BuilderProfileResponse{
			ID:          builderProfile.ID,
			CompanyName: builderProfile.Company.Name,
			DisplayName: builderProfile.DisplayName,
			Location:    builderProfile.Location,
			Phone:       nil, // Not available in BuilderProfile model
			Email:       nil, // Not available in BuilderProfile model
			CreatedAt:   builderProfile.CreatedAt,
			UpdatedAt:   builderProfile.UpdatedAt,
		}
	}

	// Add jobsite information
	if jobsite != nil {
		jobResp.Jobsite = &payload.JobsiteResponse{
			ID:          jobsite.ID,
			Name:        getStringValue(jobsite.Description), // Use description as name
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
		log.Printf("🔍 convertToJobResponseWithRelations - Jobsite added: %+v", jobResp.Jobsite)
	} else {
		log.Printf("🚫 convertToJobResponseWithRelations - Jobsite is nil")
	}

	// Add job type information
	if jobType != nil {
		jobResp.JobType = &payload.JobTypeResponse{
			ID:          jobType.ID,
			Name:        jobType.Name,
			Description: jobType.Description,
			CreatedAt:   jobType.CreatedAt,
			UpdatedAt:   jobType.UpdatedAt,
		}
		log.Printf("🔍 convertToJobResponseWithRelations - Job Type added: %+v", jobResp.JobType)
	} else {
		log.Printf("🚫 convertToJobResponseWithRelations - Job Type is nil")
	}

	// Convert job licenses with full details
	for _, jobLicense := range job.JobLicenses {
		// Get license details
		licenseDetails, err := u.licenseRepo.GetByID(ctx, jobLicense.LicenseID)
		if err != nil {
			log.Printf("🚫 Failed to get license details for ID %s: %v", jobLicense.LicenseID, err)
			continue
		}

		jobLicenseResp := payload.JobLicenseResponse{
			ID:        jobLicense.ID,
			JobID:     jobLicense.JobID,
			LicenseID: jobLicense.LicenseID,
			License: &payload.LicenseResponse{
				ID:          licenseDetails.ID,
				Name:        licenseDetails.Name,
				Description: &licenseDetails.Description,
			},
			CreatedAt: jobLicense.CreatedAt,
		}
		jobResp.JobLicenses = append(jobResp.JobLicenses, jobLicenseResp)
	}

	// Convert job skills with full details
	for _, jobSkill := range job.JobSkills {
		jobSkillResp := payload.JobSkillResponse{
			ID:                 jobSkill.ID,
			JobID:              jobSkill.JobID,
			SkillCategoryID:    jobSkill.SkillCategoryID,
			SkillSubcategoryID: jobSkill.SkillSubcategoryID,
			CreatedAt:          jobSkill.CreatedAt,
		}

		// Get skill category details if available
		if jobSkill.SkillCategoryID != nil {
			skillCategory, err := u.validator.skillCategoryRepo.GetByID(ctx, *jobSkill.SkillCategoryID)
			if err != nil {
				log.Printf("🚫 Failed to get skill category details for ID %s: %v", *jobSkill.SkillCategoryID, err)
			} else {
				jobSkillResp.SkillCategory = &payload.SkillCategoryResponse{
					ID:          skillCategory.ID,
					Name:        skillCategory.Name,
					Description: &skillCategory.Description,
				}
			}
		}

		// Get skill subcategory details if available
		if jobSkill.SkillSubcategoryID != nil {
			skillSubcategory, err := u.validator.skillSubcategoryRepo.GetByID(ctx, *jobSkill.SkillSubcategoryID)
			if err != nil {
				log.Printf("🚫 Failed to get skill subcategory details for ID %s: %v", *jobSkill.SkillSubcategoryID, err)
			} else {
				jobSkillResp.SkillSubcategory = &payload.SkillSubcategoryResponse{
					ID:          skillSubcategory.ID,
					Name:        skillSubcategory.Name,
					Description: &skillSubcategory.Description,
				}
			}
		}

		jobResp.JobSkills = append(jobResp.JobSkills, jobSkillResp)
	}

	// Convert job requirements with full details
	for _, jobRequirement := range job.JobRequirements {
		// Get job requirement details
		requirementDetails, err := u.jobRequirementRepo.GetByID(ctx, jobRequirement.JobRequirementID)
		if err != nil {
			log.Printf("🚫 Failed to get job requirement details for ID %s: %v", jobRequirement.JobRequirementID, err)
			continue
		}

		jobRequirementResp := payload.JobRequirementResponse{
			ID:          requirementDetails.ID,
			Name:        requirementDetails.Name,
			Description: requirementDetails.Description,
			IsActive:    requirementDetails.IsActive,
			CreatedAt:   requirementDetails.CreatedAt,
			UpdatedAt:   requirementDetails.UpdatedAt,
		}
		jobResp.JobRequirements = append(jobResp.JobRequirements, jobRequirementResp)
	}

	return jobResp
}

// convertToJobResponse converts a Job model to JobResponse
func (u *jobUsecase) convertToJobResponse(ctx context.Context, job *models.Job) payload.JobResponse {
	jobResp := payload.JobResponse{
		ID:                          job.ID,
		ManyLabours:                 job.ManyLabours,
		OngoingWork:                 job.OngoingWork,
		WageSiteAllowance:           job.WageSiteAllowance,
		WageLeadingHandAllowance:    job.WageLeadingHandAllowance,
		WageProductivityAllowance:   job.WageProductivityAllowance,
		ExtrasOvertimeRate:          job.ExtrasOvertimeRate,
		WageHourlyRate:              job.WageHourlyRate,
		TravelAllowance:             job.TravelAllowance,
		GST:                         job.GST,
		StartDateWork:               job.StartDateWork,
		EndDateWork:                 job.EndDateWork,
		WorkSaturday:                job.WorkSaturday,
		WorkSunday:                  job.WorkSunday,
		StartTime:                   job.StartTime,
		EndTime:                     job.EndTime,
		Description:                 job.Description,
		PaymentDay:                  job.PaymentDay,
		RequiresSupervisorSignature: job.RequiresSupervisorSignature,
		SupervisorName:              job.SupervisorName,
		Visibility:                  job.Visibility,
		PaymentType:                 job.PaymentType,
		CreatedAt:                   job.CreatedAt,
		UpdatedAt:                   job.UpdatedAt,
	}

	// Convert job licenses with full details
	for _, jobLicense := range job.JobLicenses {
		// Get license details
		licenseDetails, err := u.licenseRepo.GetByID(ctx, jobLicense.LicenseID)
		if err != nil {
			log.Printf("🚫 Failed to get license details for ID %s: %v", jobLicense.LicenseID, err)
			continue
		}

		jobLicenseResp := payload.JobLicenseResponse{
			ID:        jobLicense.ID,
			JobID:     jobLicense.JobID,
			LicenseID: jobLicense.LicenseID,
			License: &payload.LicenseResponse{
				ID:          licenseDetails.ID,
				Name:        licenseDetails.Name,
				Description: &licenseDetails.Description,
			},
			CreatedAt: jobLicense.CreatedAt,
		}
		jobResp.JobLicenses = append(jobResp.JobLicenses, jobLicenseResp)
	}

	// Convert job skills (basic info only)
	for _, jobSkill := range job.JobSkills {
		jobSkillResp := payload.JobSkillResponse{
			ID:                 jobSkill.ID,
			JobID:              jobSkill.JobID,
			SkillCategoryID:    jobSkill.SkillCategoryID,
			SkillSubcategoryID: jobSkill.SkillSubcategoryID,
			CreatedAt:          jobSkill.CreatedAt,
		}
		jobResp.JobSkills = append(jobResp.JobSkills, jobSkillResp)
	}

	return jobResp
}

// GetBuilderJobDetail retrieves a job detail for a builder (only their own jobs)
func (u *jobUsecase) GetBuilderJobDetail(ctx context.Context, jobID uuid.UUID, builderProfileID uuid.UUID) (*payload.GetJobResponse, error) {
	log.Printf("🔍 GetBuilderJobDetail - Starting with Job ID: %s, Builder Profile ID: %s", jobID, builderProfileID)

	// Get job with all relations
	job, err := u.jobRepo.GetWithRelations(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found")
	}

	// Verify that the job belongs to the builder
	if job.BuilderProfileID != builderProfileID {
		return nil, fmt.Errorf("invalid job - not owned by builder")
	}

	// Get additional relation data
	builderProfile, err := u.builderRepo.GetByID(ctx, job.BuilderProfileID)
	if err != nil {
		log.Printf("🚫 GetBuilderJobDetail - Failed to get builder profile: %v", err)
		return nil, fmt.Errorf("failed to get builder profile: %w", err)
	}
	log.Printf("🔍 GetBuilderJobDetail - Builder Profile: %+v", builderProfile)

	jobsite, err := u.jobsiteRepo.GetByID(ctx, job.JobsiteID)
	if err != nil {
		log.Printf("🚫 GetBuilderJobDetail - Failed to get jobsite: %v", err)
		return nil, fmt.Errorf("failed to get jobsite: %w", err)
	}
	log.Printf("🔍 GetBuilderJobDetail - Jobsite: %+v", jobsite)

	jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID)
	if err != nil {
		log.Printf("🚫 GetBuilderJobDetail - Failed to get job type: %v", err)
		return nil, fmt.Errorf("failed to get job type: %w", err)
	}
	log.Printf("🔍 GetBuilderJobDetail - Job Type: %+v", jobType)

	// Convert job to response with additional data
	jobResp := u.convertToJobResponseWithRelations(ctx, job, builderProfile, jobsite, jobType)
	log.Printf("🔍 GetBuilderJobDetail - Job Response: %+v", jobResp)
	log.Printf("🔍 GetBuilderJobDetail - Jobsite in response: %+v", jobResp.Jobsite)
	log.Printf("🔍 GetBuilderJobDetail - JobType in response: %+v", jobResp.JobType)

	response := &payload.GetJobResponse{
		Job:     jobResp,
		Message: "Job detail retrieved successfully",
	}

	return response, nil
}

// UpdateJobVisibility updates the visibility of a job (only for the owner builder)
func (u *jobUsecase) UpdateJobVisibility(ctx context.Context, jobID uuid.UUID, builderProfileID uuid.UUID, req payload.UpdateJobVisibilityRequest) (*payload.UpdateJobVisibilityResponse, error) {
	log.Printf("🔍 UpdateJobVisibility Usecase - Job ID: %s, Builder Profile ID: %s", jobID, builderProfileID)

	// Get job to verify ownership
	job, err := u.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		log.Printf("🚫 UpdateJobVisibility Usecase - Job not found: %v", err)
		return nil, fmt.Errorf("job not found")
	}
	log.Printf("🔍 UpdateJobVisibility Usecase - Job found, Builder Profile ID: %s", job.BuilderProfileID)

	// Verify that the job belongs to the builder
	if job.BuilderProfileID != builderProfileID {
		return nil, fmt.Errorf("invalid job - not owned by builder")
	}

	// Update only the visibility field
	job.Visibility = req.Visibility
	job.UpdatedAt = time.Now()

	// Save the updated job
	if err := u.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to update job visibility: %w", err)
	}

	// Convert job to response
	jobResp := u.convertToJobResponse(ctx, job)

	response := &payload.UpdateJobVisibilityResponse{
		Job:     jobResp,
		Message: "Job visibility updated successfully",
	}

	return response, nil
}

// calculateTotalWage calculates the total wage by summing all wage components
func (u *jobUsecase) calculateTotalWage(job *models.Job) float64 {
	var total float64

	if job.WageSiteAllowance != nil {
		total += *job.WageSiteAllowance
	}
	if job.WageLeadingHandAllowance != nil {
		total += *job.WageLeadingHandAllowance
	}
	if job.WageProductivityAllowance != nil {
		total += *job.WageProductivityAllowance
	}
	if job.WageHourlyRate != nil {
		total += *job.WageHourlyRate
	}
	if job.TravelAllowance != nil {
		total += *job.TravelAllowance
	}
	if job.ExtrasOvertimeRate != nil {
		total += *job.ExtrasOvertimeRate
	}

	return total
}

// GetLabourApplicants retrieves all applications for a labour user with job information
func (u *jobUsecase) GetLabourApplicants(ctx context.Context, labourUserID uuid.UUID) (*payload.LabourApplicantsResponse, error) {
	// Get all applications for this labour user (with pagination - using large limit to get all)
	applications, _, err := u.jobApplicationRepo.GetByLabourUserID(ctx, labourUserID, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	var labourApplicants []payload.LabourApplicationInfo

	for _, application := range applications {
		// Filter only applications with status "APPLIED"
		if application.Status != job_application_models.ApplicationStatusApplied {
			continue
		}
		// Get job information for each application
		job, err := u.jobRepo.GetByID(ctx, application.JobID)
		if err != nil {
			log.Printf("🚫 Failed to get job for application %s: %v", application.ID, err)
			continue
		}

		// Build job info with basic data
		jobInfo := payload.JobInfo{
			ID:             job.ID.String(),
			Description:    getStringValue(job.Description),
			StartDate:      job.StartDateWork,
			EndDate:        job.EndDateWork,
			WageHourlyRate: job.WageHourlyRate,
			Visibility:     string(job.Visibility),
			CreatedAt:      job.CreatedAt,
		}

		// Try to get additional relations separately
		// Get builder profile info
		if builderProfile, err := u.builderRepo.GetByID(ctx, job.BuilderProfileID); err == nil {
			companyName := getCompanyName(builderProfile.Company)
			jobInfo.BuilderProfile = &payload.BuilderProfileInfo{
				ID:          builderProfile.ID.String(),
				CompanyName: &companyName,
				DisplayName: getStringValue(builderProfile.DisplayName),
				Location:    builderProfile.Location,
			}
		}

		// Get jobsite info
		if jobsite, err := u.jobsiteRepo.GetByID(ctx, job.JobsiteID); err == nil {
			jobInfo.Jobsite = &payload.JobsiteApplicationInfo{
				ID:          jobsite.ID.String(),
				Name:        jobsite.Description,
				Address:     jobsite.Address,
				City:        jobsite.City,
				Suburb:      jobsite.Suburb,
				Description: jobsite.Description,
			}
		}

		// Get job type info
		if jobType, err := u.jobTypeRepo.GetByID(ctx, job.JobTypeID); err == nil {
			jobInfo.JobType = &payload.JobTypeInfo{
				ID:          jobType.ID.String(),
				Name:        jobType.Name,
				Description: jobType.Description,
			}
		}

		// Build labour applicant info
		labourApplicant := payload.LabourApplicationInfo{
			ApplicationID: application.ID.String(),
			JobID:         application.JobID.String(),
			Status:        string(application.Status),
			CoverLetter:   application.CoverLetter,
			ExpectedRate:  application.ExpectedRate,
			ResumeURL:     application.ResumeURL,
			AppliedAt:     application.CreatedAt,
			Job:           jobInfo,
		}

		labourApplicants = append(labourApplicants, labourApplicant)
	}

	response := &payload.LabourApplicantsResponse{
		Applications: labourApplicants,
		Total:        len(labourApplicants),
		Message:      "Applications retrieved successfully",
	}

	return response, nil
}
