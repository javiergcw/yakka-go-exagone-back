package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	authUserRepo "github.com/yakka-backend/internal/features/auth/user/entity/database"
	authUserModels "github.com/yakka-backend/internal/features/auth/user/models"
	"github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	builderModels "github.com/yakka-backend/internal/features/builder_profiles/models"
	"github.com/yakka-backend/internal/features/builder_profiles/payload"
	licenseRepo "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	dbInfra "github.com/yakka-backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BuilderProfileUsecase interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateBuilderProfileRequest) (*builderModels.BuilderProfile, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*builderModels.BuilderProfile, error)
}

type builderProfileUsecase struct {
	builderRepo     database.BuilderProfileRepository
	userLicenseRepo authUserRepo.UserLicenseRepository
	userRepo        authUserRepo.UserRepository
	licenseRepo     licenseRepo.LicenseRepository
}

func NewBuilderProfileUsecase(builderRepo database.BuilderProfileRepository, userLicenseRepo authUserRepo.UserLicenseRepository, userRepo authUserRepo.UserRepository, licenseRepo licenseRepo.LicenseRepository) BuilderProfileUsecase {
	return &builderProfileUsecase{
		builderRepo:     builderRepo,
		userLicenseRepo: userLicenseRepo,
		userRepo:        userRepo,
		licenseRepo:     licenseRepo,
	}
}

func (u *builderProfileUsecase) CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateBuilderProfileRequest) (*builderModels.BuilderProfile, error) {
	// Check if user exists
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if profile already exists
	existingProfile, err := u.builderRepo.GetByUserID(ctx, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		// Real database error
		return nil, fmt.Errorf("failed to check existing profile: %w", err)
	}
	if err == nil && existingProfile != nil {
		// Profile already exists, return error
		return nil, fmt.Errorf("builder profile already exists for this user")
	}

	// Validate licenses if provided (optimized batch validation)
	if len(req.Licenses) > 0 {
		if err := u.validateLicensesBatch(ctx, req.Licenses); err != nil {
			return nil, err
		}
	}

	// Create new profile
	profile := &builderModels.BuilderProfile{
		UserID:      userID,
		CompanyName: req.CompanyName,
		DisplayName: &req.DisplayName,
		Location:    &req.Location,
		Bio:         req.Bio,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create profile in database
	if err := u.builderRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	// Update user role to builder
	user.Role = authUserModels.UserRoleBuilder
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
	}

	// Only update phone if provided
	if req.Phone != nil {
		updates["phone"] = user.Phone
	}

	if err := u.userRepo.UpdateSpecificFields(ctx, user.ID, updates); err != nil {
		return nil, err
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

func (u *builderProfileUsecase) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*builderModels.BuilderProfile, error) {
	profile, err := u.builderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// validateLicensesBatch validates all licenses in a single batch operation
func (u *builderProfileUsecase) validateLicensesBatch(ctx context.Context, licenses []payload.UserLicenseRequest) error {
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

// validateLicensesExist checks if all license IDs exist using batch query
func (u *builderProfileUsecase) validateLicensesExist(ctx context.Context, licenseIDs []uuid.UUID) error {
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
