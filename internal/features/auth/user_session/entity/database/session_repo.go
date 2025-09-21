package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user_session/models"
	"gorm.io/gorm"
)

// SessionRepository defines the interface for session data operations
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	GetByRefreshToken(ctx context.Context, tokenHash string) (*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	Revoke(ctx context.Context, id uuid.UUID) error
}

// sessionRepository implements SessionRepository
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

// Create creates a new session
func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID retrieves a session by ID
func (r *sessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByUserID retrieves all sessions for a user
func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	err := r.db.WithContext(ctx).Where("user_id = ? AND revoked_at IS NULL", userID).Find(&sessions).Error
	return sessions, err
}

// GetByRefreshToken retrieves a session by refresh token hash
func (r *sessionRepository) GetByRefreshToken(ctx context.Context, tokenHash string) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).Where("refresh_token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Update updates a session
func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete deletes a session
func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Session{}).Error
}

// DeleteByUserID deletes all sessions for a user
func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

// DeleteExpired deletes all expired sessions
func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

// Revoke revokes a session
func (r *sessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Session{}).Where("id = ?", id).Update("revoked_at", time.Now()).Error
}
