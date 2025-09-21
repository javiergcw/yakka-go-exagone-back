package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user/models"
	"github.com/yakka-backend/internal/features/labour_profiles/entity/database"
	"github.com/yakka-backend/internal/features/labour_profiles/models"
	"github.com/yakka-backend/internal/features/labour_profiles/payload"
	authUserRepo "github.com/yakka-backend/internal/features/auth/user/entity/database"
	"gorm.io/gorm"
)

type LabourProfileUsecase interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateLabourProfileRequest) (*models.LabourProfile, error)
}

type labourProfileUsecase struct {
	labourRepo database.LabourProfileRepository
	userRepo   authUserRepo.UserRepository
}

func NewLabourProfileUsecase(labourRepo database.LabourProfileRepository, userRepo authUserRepo.UserRepository) LabourProfileUsecase {
	return &labourProfileUsecase{
		labourRepo: labourRepo,
		userRepo:    userRepo,
	}
}

func (u *labourProfileUsecase) CreateProfile(ctx context.Context, userID uuid.UUID, req payload.CreateLabourProfileRequest) (*models.LabourProfile, error) {
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
	profile := &models.LabourProfile{
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Location:  req.Location,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create profile in database
	if err := u.labourRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	// Update user role to labour
	user.Role = models.UserRoleLabour
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
