package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user/entity/database"
	"github.com/yakka-backend/internal/features/auth/user/models"
	"github.com/yakka-backend/internal/features/auth/user/payload"
	builderRepo "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	labourRepo "github.com/yakka-backend/internal/features/labour_profiles/entity/database"
	"github.com/yakka-backend/internal/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthUsecase defines the interface for authentication operations
type AuthUsecase interface {
	Register(ctx context.Context, email, password string, phone *string) (*models.User, error)
	RegisterWithAutoVerify(ctx context.Context, email, password string, phone *string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	GetUserProfileInfo(ctx context.Context, userID uuid.UUID) (payload.ProfileInfo, error)
}

// authUsecase implements AuthUsecase
type authUsecase struct {
	userRepo        database.UserRepository
	builderRepo     builderRepo.BuilderProfileRepository
	labourRepo      labourRepo.LabourProfileRepository
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo database.UserRepository, builderRepo builderRepo.BuilderProfileRepository, labourRepo labourRepo.LabourProfileRepository) AuthUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		builderRepo: builderRepo,
		labourRepo:  labourRepo,
	}
}

// Register creates a new user
func (u *authUsecase) Register(ctx context.Context, email, password string, phone *string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrConflict
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Status:       models.UserStatusPending,
		Role:         models.UserRoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return user, nil
}

// RegisterWithAutoVerify creates a new user and automatically verifies the account
func (u *authUsecase) RegisterWithAutoVerify(ctx context.Context, email, password string, phone *string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrConflict
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal
	}

	// Create user with active status (verified)
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Status:       models.UserStatusActive, // Verificado automáticamente
		Role:         models.UserRoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return user, nil
}

// Login authenticates a user
func (u *authUsecase) Login(ctx context.Context, email, password string) (*models.User, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return nil, errors.ErrForbidden
	}

	// Update last login
	err = u.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		// Log error but don't fail login
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (u *authUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return user, nil
}

// UpdateUser updates a user
func (u *authUsecase) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return u.userRepo.Update(ctx, user)
}

// ChangePassword changes a user's password
func (u *authUsecase) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.ErrNotFound
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.ErrUnauthorized
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternal
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(ctx, user)
}

// UpdateLastLogin updates the last login timestamp
func (u *authUsecase) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return u.userRepo.UpdateLastLogin(ctx, userID)
}

// GetUserProfileInfo retrieves information about user profiles
func (u *authUsecase) GetUserProfileInfo(ctx context.Context, userID uuid.UUID) (payload.ProfileInfo, error) {
	profileInfo := payload.ProfileInfo{
		HasBuilderProfile: false,
		HasLabourProfile:  false,
		HasAnyProfile:     false,
	}

	// Check if user has builder profile
	_, err := u.builderRepo.GetByUserID(ctx, userID)
	if err == nil {
		profileInfo.HasBuilderProfile = true
		profileInfo.HasAnyProfile = true
	}

	// Check if user has labour profile
	_, err = u.labourRepo.GetByUserID(ctx, userID)
	if err == nil {
		profileInfo.HasLabourProfile = true
		profileInfo.HasAnyProfile = true
	}

	return profileInfo, nil
}
