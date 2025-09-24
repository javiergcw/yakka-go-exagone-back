package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	"github.com/yakka-backend/internal/features/jobs/payload"
	jobsite_db "github.com/yakka-backend/internal/features/jobsites/entity/database"
	job_type_db "github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	license_db "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	skill_category_db "github.com/yakka-backend/internal/features/masters/skills/entity/database"
)

// JobValidationService handles validation of job relationships
type JobValidationService struct {
	builderRepo          builder_db.BuilderProfileRepository
	jobsiteRepo          jobsite_db.JobsiteRepository
	jobTypeRepo          job_type_db.JobTypeRepository
	licenseRepo          license_db.LicenseRepository
	skillCategoryRepo    skill_category_db.SkillCategoryRepository
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository
}

// NewJobValidationService creates a new job validation service
func NewJobValidationService(
	builderRepo builder_db.BuilderProfileRepository,
	jobsiteRepo jobsite_db.JobsiteRepository,
	jobTypeRepo job_type_db.JobTypeRepository,
	licenseRepo license_db.LicenseRepository,
	skillCategoryRepo skill_category_db.SkillCategoryRepository,
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository,
) *JobValidationService {
	return &JobValidationService{
		builderRepo:          builderRepo,
		jobsiteRepo:          jobsiteRepo,
		jobTypeRepo:          jobTypeRepo,
		licenseRepo:          licenseRepo,
		skillCategoryRepo:    skillCategoryRepo,
		skillSubcategoryRepo: skillSubcategoryRepo,
	}
}

// ValidateCreateJobRequest validates a create job request
func (v *JobValidationService) ValidateCreateJobRequest(ctx context.Context, req payload.CreateJobRequest) error {
	// Validate builder profile exists (only if provided)
	if req.BuilderProfileID != uuid.Nil {
		if err := v.validateBuilderProfile(ctx, req.BuilderProfileID); err != nil {
			return fmt.Errorf("invalid builder profile: %w", err)
		}
	}

	// Validate jobsite exists and belongs to builder
	if req.BuilderProfileID != uuid.Nil {
		if err := v.validateJobsite(ctx, req.JobsiteID, req.BuilderProfileID); err != nil {
			return fmt.Errorf("invalid jobsite: %w", err)
		}
	}

	// Validate job type exists
	if err := v.validateJobType(ctx, req.JobTypeID); err != nil {
		return fmt.Errorf("invalid job type: %w", err)
	}

	// Validate licenses exist
	for _, licenseID := range req.LicenseIDs {
		if err := v.validateLicense(ctx, licenseID); err != nil {
			return fmt.Errorf("invalid license ID %s: %w", licenseID, err)
		}
	}

	// Validate skill categories exist
	for _, skillCategoryID := range req.SkillCategoryIDs {
		if err := v.validateSkillCategory(ctx, skillCategoryID); err != nil {
			return fmt.Errorf("invalid skill category ID %s: %w", skillCategoryID, err)
		}
	}

	// Validate skill subcategories exist
	for _, skillSubcategoryID := range req.SkillSubcategoryIDs {
		if err := v.validateSkillSubcategory(ctx, skillSubcategoryID); err != nil {
			return fmt.Errorf("invalid skill subcategory ID %s: %w", skillSubcategoryID, err)
		}
	}

	// Validate business rules
	if err := v.validateBusinessRules(req); err != nil {
		return fmt.Errorf("business rule validation failed: %w", err)
	}

	return nil
}

// validateBuilderProfile checks if builder profile exists
func (v *JobValidationService) validateBuilderProfile(ctx context.Context, builderProfileID uuid.UUID) error {
	if builderProfileID == uuid.Nil {
		return fmt.Errorf("builder profile ID cannot be nil")
	}

	// Check if builder profile exists in database
	_, err := v.builderRepo.GetByID(ctx, builderProfileID)
	if err != nil {
		return fmt.Errorf("builder profile with ID %s does not exist: %w", builderProfileID, err)
	}

	return nil
}

// validateJobsite checks if jobsite exists and belongs to the builder
func (v *JobValidationService) validateJobsite(ctx context.Context, jobsiteID uuid.UUID, builderProfileID uuid.UUID) error {
	if jobsiteID == uuid.Nil {
		return fmt.Errorf("jobsite ID cannot be nil")
	}

	// Check if jobsite exists in database
	jobsite, err := v.jobsiteRepo.GetByID(ctx, jobsiteID)
	if err != nil {
		return fmt.Errorf("jobsite with ID %s does not exist: %w", jobsiteID, err)
	}

	// Check if jobsite belongs to the builder
	if jobsite.BuilderID != builderProfileID {
		return fmt.Errorf("jobsite with ID %s does not belong to builder %s", jobsiteID, builderProfileID)
	}

	return nil
}

// validateJobType checks if job type exists
func (v *JobValidationService) validateJobType(ctx context.Context, jobTypeID uuid.UUID) error {
	if jobTypeID == uuid.Nil {
		return fmt.Errorf("job type ID cannot be nil")
	}

	// Check if job type exists in database
	_, err := v.jobTypeRepo.GetByID(ctx, jobTypeID)
	if err != nil {
		return fmt.Errorf("job type with ID %s does not exist: %w", jobTypeID, err)
	}

	return nil
}

// validateLicense checks if license exists
func (v *JobValidationService) validateLicense(ctx context.Context, licenseID uuid.UUID) error {
	if licenseID == uuid.Nil {
		return fmt.Errorf("license ID cannot be nil")
	}

	// Check if license exists in database
	_, err := v.licenseRepo.GetByID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("license with ID %s does not exist: %w", licenseID, err)
	}

	return nil
}

// validateSkillCategory checks if skill category exists
func (v *JobValidationService) validateSkillCategory(ctx context.Context, skillCategoryID uuid.UUID) error {
	if skillCategoryID == uuid.Nil {
		return fmt.Errorf("skill category ID cannot be nil")
	}

	// Check if skill category exists in database
	_, err := v.skillCategoryRepo.GetByID(ctx, skillCategoryID)
	if err != nil {
		return fmt.Errorf("skill category with ID %s does not exist: %w", skillCategoryID, err)
	}

	return nil
}

// validateSkillSubcategory checks if skill subcategory exists
func (v *JobValidationService) validateSkillSubcategory(ctx context.Context, skillSubcategoryID uuid.UUID) error {
	if skillSubcategoryID == uuid.Nil {
		return fmt.Errorf("skill subcategory ID cannot be nil")
	}

	// Check if skill subcategory exists in database
	_, err := v.skillSubcategoryRepo.GetByID(ctx, skillSubcategoryID)
	if err != nil {
		return fmt.Errorf("skill subcategory with ID %s does not exist: %w", skillSubcategoryID, err)
	}

	return nil
}

// validateBusinessRules validates business logic rules
func (v *JobValidationService) validateBusinessRules(req payload.CreateJobRequest) error {
	// Validate many labours is positive
	if req.ManyLabours <= 0 {
		return fmt.Errorf("many labours must be greater than 0")
	}

	// Validate date range if both dates are provided
	if req.StartDateWork != nil && req.EndDateWork != nil {
		if req.StartDateWork.After(*req.EndDateWork) {
			return fmt.Errorf("start date cannot be after end date")
		}
	}

	// Validate payment day for FIXED_DAY payment type
	if req.PaymentType == "FIXED_DAY" {
		if req.PaymentDay == nil || *req.PaymentDay < 1 || *req.PaymentDay > 31 {
			return fmt.Errorf("payment day must be between 1 and 31 for FIXED_DAY payment type")
		}
	}

	// Validate time format if provided
	if req.StartTime != nil {
		if err := v.validateTimeFormat(*req.StartTime); err != nil {
			return fmt.Errorf("invalid start time format: %w", err)
		}
	}

	if req.EndTime != nil {
		if err := v.validateTimeFormat(*req.EndTime); err != nil {
			return fmt.Errorf("invalid end time format: %w", err)
		}
	}

	return nil
}

// validateTimeFormat validates time format (HH:MM:SS)
func (v *JobValidationService) validateTimeFormat(timeStr string) error {
	// Basic validation - should be in HH:MM:SS format
	if len(timeStr) != 8 {
		return fmt.Errorf("time must be in HH:MM:SS format")
	}
	// TODO: Add more comprehensive time validation
	return nil
}
