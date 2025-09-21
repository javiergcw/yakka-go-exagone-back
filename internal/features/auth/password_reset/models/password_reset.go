package models

import (
	"time"

	"github.com/google/uuid"
)

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"not null;type:uuid"`
	TokenHash string     `json:"-" gorm:"not null;type:text"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;type:timestamptz"`
	UsedAt    *time.Time `json:"used_at" gorm:"type:timestamptz"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the PasswordReset model
func (PasswordReset) TableName() string {
	return "password_resets"
}
