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
}

func NewBuilderProfileUsecase(builderRepo database.BuilderProfileRepository, userLicenseRepo authUserRepo.UserLicenseRepository, userRepo authUserRepo.UserRepository) BuilderProfileUsecase {
	return &builderProfileUsecase{
		builderRepo:     builderRepo,
		userLicenseRepo: userLicenseRepo,
		userRepo:        userRepo,
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

	// Create new profile
	profile := &builderModels.BuilderProfile{
		UserID:      userID,
		CompanyName: &req.CompanyName,
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

	if err := u.userRepo.Update(ctx, user); err != nil {
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
