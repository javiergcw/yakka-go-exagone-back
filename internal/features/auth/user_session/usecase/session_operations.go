package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user_session/entity/database"
	"github.com/yakka-backend/internal/features/auth/user_session/models"
	"github.com/yakka-backend/internal/shared/errors"
)

// SessionUsecase defines the interface for session operations
type SessionUsecase interface {
	CreateSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*models.Session, string, error)
	RefreshSession(ctx context.Context, refreshToken string) (*models.Session, string, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	CleanupExpiredSessions(ctx context.Context) error
}

// sessionUsecase implements SessionUsecase
type sessionUsecase struct {
	sessionRepo database.SessionRepository
}

// NewSessionUsecase creates a new session usecase
func NewSessionUsecase(sessionRepo database.SessionRepository) SessionUsecase {
	return &sessionUsecase{
		sessionRepo: sessionRepo,
	}
}

// CreateSession creates a new session for a user
func (u *sessionUsecase) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*models.Session, string, error) {
	// Generate refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	// Hash the token for storage
	tokenHash := hashToken(refreshToken)

	// Create session
	session := &models.Session{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // 7 days
		UserAgent:        &userAgent,
		IPAddress:        &ipAddress,
		CreatedAt:        time.Now(),
	}

	err = u.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	return session, refreshToken, nil
}

// RefreshSession refreshes a session using a refresh token
func (u *sessionUsecase) RefreshSession(ctx context.Context, refreshToken string) (*models.Session, string, error) {
	// Hash the token to find the session
	tokenHash := hashToken(refreshToken)

	// Get session by token hash
	session, err := u.sessionRepo.GetByRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, "", errors.ErrUnauthorized
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	// Update session with new token
	session.RefreshTokenHash = hashToken(newRefreshToken)
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // 7 days

	err = u.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	return session, newRefreshToken, nil
}

// GetSession retrieves a session by ID
func (u *sessionUsecase) GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	session, err := u.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return session, nil
}

// RevokeSession revokes a specific session
func (u *sessionUsecase) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return u.sessionRepo.Revoke(ctx, sessionID)
}

// RevokeAllUserSessions revokes all sessions for a user
func (u *sessionUsecase) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return u.sessionRepo.DeleteByUserID(ctx, userID)
}

// CleanupExpiredSessions removes all expired sessions
func (u *sessionUsecase) CleanupExpiredSessions(ctx context.Context) error {
	return u.sessionRepo.DeleteExpired(ctx)
}

// generateRefreshToken generates a secure random refresh token
func generateRefreshToken() (string, error) {
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
