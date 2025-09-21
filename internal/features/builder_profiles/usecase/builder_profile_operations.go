package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user/models"
	"github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	"github.com/yakka-backend/internal/features/builder_profiles/models"
	"github.com/yakka-backend/internal/features/builder_profiles/payload"
	authUserRepo "github.com/yakka-backend/internal/features/auth/user/entity/database"
	"gorm.io/gorm"
)

type BuilderProfileUsecase interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateBuilderProfileRequest) (*models.BuilderProfile, error)
}

type builderProfileUsecase struct {
	builderRepo database.BuilderProfileRepository
	userRepo    authUserRepo.UserRepository
}

func NewBuilderProfileUsecase(builderRepo database.BuilderProfileRepository, userRepo authUserRepo.UserRepository) BuilderProfileUsecase {
	return &builderProfileUsecase{
		builderRepo: builderRepo,
		userRepo:     userRepo,
	}
}

func (u *builderProfileUsecase) CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateBuilderProfileRequest) (*models.BuilderProfile, error) {
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
	profile := &models.BuilderProfile{
		UserID:      userID,
		CompanyName: req.CompanyName,
		DisplayName: req.DisplayName,
		Location:    req.Location,
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create profile in database
	if err := u.builderRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	// Update user role to builder
	user.Role = models.UserRoleBuilder
	user.RoleChangedAt = &time.Time{}
	*user.RoleChangedAt = time.Now()

	// Update user phone if provided
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return profile, nil
}
