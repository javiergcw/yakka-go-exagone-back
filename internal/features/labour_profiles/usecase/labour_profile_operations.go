package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	authUserRepo "github.com/yakka-backend/internal/features/auth/user/entity/database"
	authUserModels "github.com/yakka-backend/internal/features/auth/user/models"
	"github.com/yakka-backend/internal/features/labour_profiles/entity/database"
	labourModels "github.com/yakka-backend/internal/features/labour_profiles/models"
	"github.com/yakka-backend/internal/features/labour_profiles/payload"
	experienceRepo "github.com/yakka-backend/internal/features/masters/experience_levels/entity/database"
	licenseRepo "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	skillRepo "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	dbInfra "github.com/yakka-backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type LabourProfileUsecase interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateLabourProfileRequest) (*labourModels.LabourProfile, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*labourModels.LabourProfile, error)
}

type labourProfileUsecase struct {
	labourRepo           database.LabourProfileRepository
	labourSkillRepo      database.LabourProfileSkillRepository
	userLicenseRepo      authUserRepo.UserLicenseRepository
	userRepo             authUserRepo.UserRepository
	licenseRepo          licenseRepo.LicenseRepository
	skillCategoryRepo    skillRepo.SkillCategoryRepository
	skillSubcategoryRepo skillRepo.SkillSubcategoryRepository
	experienceRepo       experienceRepo.ExperienceLevelRepository
}

func NewLabourProfileUsecase(labourRepo database.LabourProfileRepository, labourSkillRepo database.LabourProfileSkillRepository, userLicenseRepo authUserRepo.UserLicenseRepository, userRepo authUserRepo.UserRepository, licenseRepo licenseRepo.LicenseRepository, skillCategoryRepo skillRepo.SkillCategoryRepository, skillSubcategoryRepo skillRepo.SkillSubcategoryRepository, experienceRepo experienceRepo.ExperienceLevelRepository) LabourProfileUsecase {
	return &labourProfileUsecase{
		labourRepo:           labourRepo,
		labourSkillRepo:      labourSkillRepo,
		userLicenseRepo:      userLicenseRepo,
		userRepo:             userRepo,
		licenseRepo:          licenseRepo,
		skillCategoryRepo:    skillCategoryRepo,
		skillSubcategoryRepo: skillSubcategoryRepo,
		experienceRepo:       experienceRepo,
	}
}

func (u *labourProfileUsecase) CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateLabourProfileRequest) (*labourModels.LabourProfile, error) {
	// Check if user exists
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if profile already exists
	existingProfile, err := u.labourRepo.GetByUserID(ctx, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		// Real database error
		return nil, fmt.Errorf("failed to check existing profile: %w", err)
	}
	if err == nil && existingProfile != nil {
		// Profile already exists, return error
		return nil, fmt.Errorf("labour profile already exists for this user")
	}

	// Validate skills if provided (optimized batch validation)
	if len(req.Skills) > 0 {
		if err := u.validateSkillsBatch(ctx, req.Skills); err != nil {
			return nil, err
		}
	}

	// Validate licenses if provided (optimized batch validation)
	if len(req.Licenses) > 0 {
		if err := u.validateLicensesBatch(ctx, req.Licenses); err != nil {
			return nil, err
		}
	}

	// Create new profile
	profile := &labourModels.LabourProfile{
		UserID:    userID,
		Location:  &req.Location,
		Bio:       req.Bio,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create profile in database
	if err := u.labourRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	// Update user role to labour
	user.Role = authUserModels.UserRoleLabour
	user.RoleChangedAt = &time.Time{}
	*user.RoleChangedAt = time.Now()

	// Update user fields that are now in the user table
	user.Photo = req.AvatarURL

	// Update user phone if provided
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	// Use Updates instead of Save to avoid overwriting existing fields
	updates := map[string]interface{}{
		"role":            user.Role,
		"role_changed_at": user.RoleChangedAt,
		"photo":           user.Photo,
		"first_name":      req.FirstName,
		"last_name":       req.LastName,
	}

	// Only update phone if provided
	if req.Phone != nil {
		updates["phone"] = user.Phone
	}

	if err := u.userRepo.UpdateSpecificFields(ctx, user.ID, updates); err != nil {
		return nil, err
	}

	// Create skills if provided
	if len(req.Skills) > 0 {
		var skills []*labourModels.LabourProfileSkill
		for _, skillReq := range req.Skills {
			categoryID, err := uuid.Parse(skillReq.CategoryID)
			if err != nil {
				return nil, fmt.Errorf("invalid category_id: %w", err)
			}

			subcategoryID, err := uuid.Parse(skillReq.SubcategoryID)
			if err != nil {
				return nil, fmt.Errorf("invalid subcategory_id: %w", err)
			}

			experienceLevelID, err := uuid.Parse(skillReq.ExperienceLevelID)
			if err != nil {
				return nil, fmt.Errorf("invalid experience_level_id: %w", err)
			}

			// Validate that category exists
			_, err = u.skillCategoryRepo.GetByID(ctx, categoryID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, fmt.Errorf("skill category not found")
				}
				return nil, fmt.Errorf("failed to validate skill category: %w", err)
			}

			// Validate that subcategory exists
			_, err = u.skillSubcategoryRepo.GetByID(ctx, subcategoryID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, fmt.Errorf("skill subcategory not found")
				}
				return nil, fmt.Errorf("failed to validate skill subcategory: %w", err)
			}

			// Validate that experience level exists
			_, err = u.experienceRepo.GetByID(ctx, experienceLevelID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, fmt.Errorf("experience level not found")
				}
				return nil, fmt.Errorf("failed to validate experience level: %w", err)
			}

			skill := &labourModels.LabourProfileSkill{
				LabourProfileID:   profile.ID,
				CategoryID:        categoryID,
				SubcategoryID:     subcategoryID,
				ExperienceLevelID: experienceLevelID,
				YearsExperience:   skillReq.YearsExperience,
				IsPrimary:         skillReq.IsPrimary,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			skills = append(skills, skill)
		}

		if err := u.labourSkillRepo.CreateBatch(ctx, skills); err != nil {
			// If skills creation fails, we should rollback the profile creation
			// For now, we'll log the error but not fail the entire operation
			// In a production system, you might want to use database transactions
			fmt.Printf("Warning: Failed to create skills for profile %s: %v\n", profile.ID, err)
		}
	}

	// Create licenses if provided
	if len(req.Licenses) > 0 {
		var licenses []*authUserModels.UserLicense
		for _, licenseReq := range req.Licenses {
			licenseID, err := uuid.Parse(licenseReq.LicenseID)
			if err != nil {
				return nil, fmt.Errorf("invalid license_id: %w", err)
			}

			// Validate that license exists
			_, err = u.licenseRepo.GetByID(ctx, licenseID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, fmt.Errorf("license not found")
				}
				return nil, fmt.Errorf("failed to validate license: %w", err)
			}

			license := &authUserModels.UserLicense{
				UserID:    userID,
				LicenseID: licenseID,
				PhotoURL:  licenseReq.PhotoURL,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Parse dates if provided
			if licenseReq.IssuedAt != nil {
				if issuedAt, err := time.Parse("2006-01-02T15:04:05Z07:00", *licenseReq.IssuedAt); err == nil {
					license.IssuedAt = &issuedAt
				}
			}
			if licenseReq.ExpiresAt != nil {
				if expiresAt, err := time.Parse("2006-01-02T15:04:05Z07:00", *licenseReq.ExpiresAt); err == nil {
					license.ExpiresAt = &expiresAt
				}
			}

			licenses = append(licenses, license)
		}

		if err := u.userLicenseRepo.CreateBatch(ctx, licenses); err != nil {
			// If licenses creation fails, we should rollback the profile creation
			// For now, we'll log the error but not fail the entire operation
			// In a production system, you might want to use database transactions
			fmt.Printf("Warning: Failed to create licenses for user %s: %v\n", userID, err)
		}
	}

	return profile, nil
}

func (u *labourProfileUsecase) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*labourModels.LabourProfile, error) {
	profile, err := u.labourRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// validateSkillsBatch validates all skills in a single batch operation
func (u *labourProfileUsecase) validateSkillsBatch(ctx context.Context, skills []payload.LabourProfileSkillRequest) error {
	// Collect all unique IDs
	categoryIDs := make([]uuid.UUID, 0, len(skills))
	subcategoryIDs := make([]uuid.UUID, 0, len(skills))
	experienceLevelIDs := make([]uuid.UUID, 0, len(skills))

	for _, skill := range skills {
		categoryID, err := uuid.Parse(skill.CategoryID)
		if err != nil {
			return fmt.Errorf("invalid category_id: %w", err)
		}
		categoryIDs = append(categoryIDs, categoryID)

		subcategoryID, err := uuid.Parse(skill.SubcategoryID)
		if err != nil {
			return fmt.Errorf("invalid subcategory_id: %w", err)
		}
		subcategoryIDs = append(subcategoryIDs, subcategoryID)

		experienceLevelID, err := uuid.Parse(skill.ExperienceLevelID)
		if err != nil {
			return fmt.Errorf("invalid experience_level_id: %w", err)
		}
		experienceLevelIDs = append(experienceLevelIDs, experienceLevelID)
	}

	// Validate categories exist (single query)
	if err := u.validateCategoriesExist(ctx, categoryIDs); err != nil {
		return err
	}

	// Validate subcategories exist (single query)
	if err := u.validateSubcategoriesExist(ctx, subcategoryIDs); err != nil {
		return err
	}

	// Validate experience levels exist (single query)
	if err := u.validateExperienceLevelsExist(ctx, experienceLevelIDs); err != nil {
		return err
	}

	return nil
}

// validateLicensesBatch validates all licenses in a single batch operation
func (u *labourProfileUsecase) validateLicensesBatch(ctx context.Context, licenses []payload.UserLicenseRequest) error {
	licenseIDs := make([]uuid.UUID, 0, len(licenses))

	for _, license := range licenses {
		licenseID, err := uuid.Parse(license.LicenseID)
		if err != nil {
			return fmt.Errorf("invalid license_id: %w", err)
		}
		licenseIDs = append(licenseIDs, licenseID)
	}

	// Validate licenses exist (single query)
	if err := u.validateLicensesExist(ctx, licenseIDs); err != nil {
		return err
	}

	return nil
}

// validateCategoriesExist checks if all category IDs exist using batch query
func (u *labourProfileUsecase) validateCategoriesExist(ctx context.Context, categoryIDs []uuid.UUID) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	// Use raw SQL with optimized query for maximum performance
	var count int64
	err := dbInfra.DB.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM skill_categories WHERE id IN (?) AND deleted_at IS NULL", categoryIDs).
		Scan(&count).Error

	if err != nil {
		return fmt.Errorf("failed to validate skill categories: %w", err)
	}

	if int(count) != len(categoryIDs) {
		return fmt.Errorf("skill category not found")
	}

	return nil
}

// validateSubcategoriesExist checks if all subcategory IDs exist using batch query
func (u *labourProfileUsecase) validateSubcategoriesExist(ctx context.Context, subcategoryIDs []uuid.UUID) error {
	if len(subcategoryIDs) == 0 {
		return nil
	}

	// Use raw SQL with optimized query for maximum performance
	var count int64
	err := dbInfra.DB.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM skill_subcategories WHERE id IN (?) AND deleted_at IS NULL", subcategoryIDs).
		Scan(&count).Error

	if err != nil {
		return fmt.Errorf("failed to validate skill subcategories: %w", err)
	}

	if int(count) != len(subcategoryIDs) {
		return fmt.Errorf("skill subcategory not found")
	}

	return nil
}

// validateExperienceLevelsExist checks if all experience level IDs exist using batch query
func (u *labourProfileUsecase) validateExperienceLevelsExist(ctx context.Context, experienceLevelIDs []uuid.UUID) error {
	if len(experienceLevelIDs) == 0 {
		return nil
	}

	// Use raw SQL with optimized query for maximum performance
	var count int64
	err := dbInfra.DB.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM experience_levels WHERE id IN (?) AND deleted_at IS NULL", experienceLevelIDs).
		Scan(&count).Error

	if err != nil {
		return fmt.Errorf("failed to validate experience levels: %w", err)
	}

	if int(count) != len(experienceLevelIDs) {
		return fmt.Errorf("experience level not found")
	}

	return nil
}

// validateLicensesExist checks if all license IDs exist using batch query
func (u *labourProfileUsecase) validateLicensesExist(ctx context.Context, licenseIDs []uuid.UUID) error {
	if len(licenseIDs) == 0 {
		return nil
	}

	// Use raw SQL for optimal performance
	var count int64
	err := dbInfra.DB.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM licenses WHERE id IN (?)", licenseIDs).
		Scan(&count).Error

	if err != nil {
		return fmt.Errorf("failed to validate licenses: %w", err)
	}

	if int(count) != len(licenseIDs) {
		return fmt.Errorf("license not found")
	}

	return nil
}
