package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/password_reset/models"
	"gorm.io/gorm"
)

// PasswordResetRepository defines the interface for password reset data operations
type PasswordResetRepository interface {
	Create(ctx context.Context, reset *models.PasswordReset) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PasswordReset, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.PasswordReset, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PasswordReset, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// passwordResetRepository implements PasswordResetRepository
type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// Create creates a new password reset request
func (r *passwordResetRepository) Create(ctx context.Context, reset *models.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

// GetByID retrieves a password reset by ID
func (r *passwordResetRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// GetByTokenHash retrieves a password reset by token hash
func (r *passwordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.WithContext(ctx).Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now()).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// GetByUserID retrieves all password resets for a user
func (r *passwordResetRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PasswordReset, error) {
	var resets []*models.PasswordReset
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&resets).Error
	return resets, err
}

// MarkAsUsed marks a password reset as used
func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.PasswordReset{}).Where("id = ?", id).Update("used_at", time.Now()).Error
}

// DeleteExpired deletes all expired password resets
func (r *passwordResetRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{}).Error
}

// DeleteByUserID deletes all password resets for a user
func (r *passwordResetRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.PasswordReset{}).Error
}
