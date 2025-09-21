package models

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type Session struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID  `json:"user_id" gorm:"not null;type:uuid;index"`
	RefreshTokenHash string     `json:"-" gorm:"not null;type:text"`
	ExpiresAt        time.Time  `json:"expires_at" gorm:"not null;type:timestamptz;index"`
	UserAgent        *string    `json:"user_agent" gorm:"size:255"`
	IPAddress        *string    `json:"ip_address" gorm:"size:45"`
	RevokedAt        *time.Time `json:"revoked_at" gorm:"type:timestamptz"`
	CreatedAt        time.Time  `json:"created_at" gorm:"not null;type:timestamptz;index"`
}

// TableName returns the table name for the Session model
func (Session) TableName() string {
	return "sessions"
}
