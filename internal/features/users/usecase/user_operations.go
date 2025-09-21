package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yakka-backend/internal/features/users"
	"github.com/yakka-backend/internal/features/users/models"
)

// userUseCase implements the UseCase interface
type userUseCase struct {
	userRepo  users.Repository
	validator *validator.Validate
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo users.Repository) *userUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// CreateUser creates a new user
func (u *userUseCase) CreateUser(ctx context.Context, user *models.User) error {
	// Validate user data
	if err := u.validator.Struct(user); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Check if email already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Create user
	if err := u.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID
func (u *userUseCase) GetUser(ctx context.Context, id uint) (*models.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetAllUsers retrieves all users
func (u *userUseCase) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

// UpdateUser updates an existing user
func (u *userUseCase) UpdateUser(ctx context.Context, user *models.User) error {
	// Validate user data
	if err := u.validator.Struct(user); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Check if user exists
	existingUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if email is being changed and if new email already exists
	if existingUser.Email != user.Email {
		emailUser, err := u.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && emailUser != nil {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	// Update user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (u *userUseCase) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete user
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
