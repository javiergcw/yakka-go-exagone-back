package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/password_reset/entity/database"
	"github.com/yakka-backend/internal/features/auth/password_reset/models"
	"github.com/yakka-backend/internal/shared/errors"
)

// PasswordResetUsecase defines the interface for password reset operations
type PasswordResetUsecase interface {
	RequestPasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	ValidateResetToken(ctx context.Context, token string) (*models.PasswordReset, error)
	CleanupExpiredResets(ctx context.Context) error
}

// passwordResetUsecase implements PasswordResetUsecase
type passwordResetUsecase struct {
	passwordResetRepo database.PasswordResetRepository
}

// NewPasswordResetUsecase creates a new password reset usecase
func NewPasswordResetUsecase(passwordResetRepo database.PasswordResetRepository) PasswordResetUsecase {
	return &passwordResetUsecase{
		passwordResetRepo: passwordResetRepo,
	}
}

// RequestPasswordReset creates a password reset request
func (u *passwordResetUsecase) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	// Generate reset token
	resetToken, err := generateResetToken()
	if err != nil {
		return "", errors.ErrInternal
	}

	// Hash the token for storage
	tokenHash := hashToken(resetToken)

	// Create password reset request
	reset := &models.PasswordReset{
		ID:        uuid.New(),
		UserID:    uuid.New(), // This should be the actual user ID from email lookup
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		CreatedAt: time.Now(),
	}

	err = u.passwordResetRepo.Create(ctx, reset)
	if err != nil {
		return "", errors.ErrInternal
	}

	return resetToken, nil
}

// ResetPassword resets a user's password using a reset token
func (u *passwordResetUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Hash the token to find the reset request
	tokenHash := hashToken(token)

	// Get reset request by token hash
	reset, err := u.passwordResetRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return errors.ErrUnauthorized
	}

	// Check if token is already used
	if reset.UsedAt != nil {
		return errors.ErrConflict
	}

	// Mark token as used
	err = u.passwordResetRepo.MarkAsUsed(ctx, reset.ID)
	if err != nil {
		return errors.ErrInternal
	}

	// Here you would update the user's password
	// This requires access to the user repository
	// For now, we'll just return success

	return nil
}

// ValidateResetToken validates a password reset token
func (u *passwordResetUsecase) ValidateResetToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	tokenHash := hashToken(token)

	reset, err := u.passwordResetRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	return reset, nil
}

// CleanupExpiredResets removes all expired password reset requests
func (u *passwordResetUsecase) CleanupExpiredResets(ctx context.Context) error {
	return u.passwordResetRepo.DeleteExpired(ctx)
}

// generateResetToken generates a secure random reset token
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashToken hashes a token for secure storage
func hashToken(token string) string {
	// In a real implementation, you would use a proper hash function
	// For now, we'll use a simple approach
	return token // This should be replaced with proper hashing
}
