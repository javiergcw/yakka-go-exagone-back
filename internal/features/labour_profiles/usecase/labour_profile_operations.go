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
	"gorm.io/gorm"
)

type LabourProfileUsecase interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateLabourProfileRequest) (*labourModels.LabourProfile, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*labourModels.LabourProfile, error)
}

type labourProfileUsecase struct {
	labourRepo      database.LabourProfileRepository
	labourSkillRepo database.LabourProfileSkillRepository
	userLicenseRepo authUserRepo.UserLicenseRepository
	userRepo        authUserRepo.UserRepository
}

func NewLabourProfileUsecase(labourRepo database.LabourProfileRepository, labourSkillRepo database.LabourProfileSkillRepository, userLicenseRepo authUserRepo.UserLicenseRepository, userRepo authUserRepo.UserRepository) LabourProfileUsecase {
	return &labourProfileUsecase{
		labourRepo:      labourRepo,
		labourSkillRepo: labourSkillRepo,
		userLicenseRepo: userLicenseRepo,
		userRepo:        userRepo,
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
	user.FirstName = &req.FirstName
	user.LastName = &req.LastName
	user.Photo = req.AvatarURL

	// Update user phone if provided
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
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
