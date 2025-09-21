package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/email_verification/models"
	"gorm.io/gorm"
)

// EmailVerificationRepository defines the interface for email verification data operations
type EmailVerificationRepository interface {
	Create(ctx context.Context, verification *models.EmailVerification) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.EmailVerification, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.EmailVerification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.EmailVerification, error)
	MarkAsVerified(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// emailVerificationRepository implements EmailVerificationRepository
type emailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository creates a new email verification repository
func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

// Create creates a new email verification request
func (r *emailVerificationRepository) Create(ctx context.Context, verification *models.EmailVerification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

// GetByID retrieves an email verification by ID
func (r *emailVerificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByTokenHash retrieves an email verification by token hash
func (r *emailVerificationRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.WithContext(ctx).Where("token_hash = ? AND verified_at IS NULL AND expires_at > ?", tokenHash, time.Now()).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByUserID retrieves all email verifications for a user
func (r *emailVerificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.EmailVerification, error) {
	var verifications []*models.EmailVerification
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&verifications).Error
	return verifications, err
}

// MarkAsVerified marks an email verification as verified
func (r *emailVerificationRepository) MarkAsVerified(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.EmailVerification{}).Where("id = ?", id).Update("verified_at", time.Now()).Error
}

// DeleteExpired deletes all expired email verifications
func (r *emailVerificationRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.EmailVerification{}).Error
}

// DeleteByUserID deletes all email verifications for a user
func (r *emailVerificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.EmailVerification{}).Error
}
