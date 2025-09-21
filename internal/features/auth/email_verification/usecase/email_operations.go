package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/email_verification/entity/database"
	"github.com/yakka-backend/internal/features/auth/email_verification/models"
	"github.com/yakka-backend/internal/shared/errors"
)

// EmailVerificationUsecase defines the interface for email verification operations
type EmailVerificationUsecase interface {
	RequestEmailVerification(ctx context.Context, userID uuid.UUID) (string, error)
	VerifyEmail(ctx context.Context, token string) error
	ValidateVerificationToken(ctx context.Context, token string) (*models.EmailVerification, error)
	CleanupExpiredVerifications(ctx context.Context) error
}

// emailVerificationUsecase implements EmailVerificationUsecase
type emailVerificationUsecase struct {
	emailVerificationRepo database.EmailVerificationRepository
	userRepo              UserRepository
}

// UserRepository interface for updating user status
type UserRepository interface {
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error
}

// NewEmailVerificationUsecase creates a new email verification usecase
func NewEmailVerificationUsecase(emailVerificationRepo database.EmailVerificationRepository, userRepo UserRepository) EmailVerificationUsecase {
	return &emailVerificationUsecase{
		emailVerificationRepo: emailVerificationRepo,
		userRepo:              userRepo,
	}
}

// RequestEmailVerification creates an email verification request
func (u *emailVerificationUsecase) RequestEmailVerification(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate verification token
	verificationToken, err := generateVerificationToken()
	if err != nil {
		return "", errors.ErrInternal
	}

	// Hash the token for storage
	tokenHash := hashToken(verificationToken)

	// Create email verification request
	verification := &models.EmailVerification{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
		CreatedAt: time.Now(),
	}

	err = u.emailVerificationRepo.Create(ctx, verification)
	if err != nil {
		return "", errors.ErrInternal
	}

	return verificationToken, nil
}

// VerifyEmail verifies a user's email using a verification token
func (u *emailVerificationUsecase) VerifyEmail(ctx context.Context, token string) error {
	// Hash the token to find the verification request
	tokenHash := hashToken(token)

	// Get verification request by token hash
	verification, err := u.emailVerificationRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return errors.ErrUnauthorized
	}

	// Check if email is already verified
	if verification.VerifiedAt != nil {
		return errors.ErrConflict
	}

	// Mark email as verified
	err = u.emailVerificationRepo.MarkAsVerified(ctx, verification.ID)
	if err != nil {
		return errors.ErrInternal
	}

	// Update user status from pending to active
	err = u.userRepo.UpdateUserStatus(ctx, verification.UserID, "active")
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

// ValidateVerificationToken validates an email verification token
func (u *emailVerificationUsecase) ValidateVerificationToken(ctx context.Context, token string) (*models.EmailVerification, error) {
	tokenHash := hashToken(token)

	verification, err := u.emailVerificationRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	return verification, nil
}

// CleanupExpiredVerifications removes all expired email verification requests
func (u *emailVerificationUsecase) CleanupExpiredVerifications(ctx context.Context) error {
	return u.emailVerificationRepo.DeleteExpired(ctx)
}

// generateVerificationToken generates a secure random verification token
func generateVerificationToken() (string, error) {
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
