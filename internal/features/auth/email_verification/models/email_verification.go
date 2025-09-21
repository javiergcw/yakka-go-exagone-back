package models

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerification represents an email verification request
type EmailVerification struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID  `json:"user_id" gorm:"not null;type:uuid"`
	TokenHash  string     `json:"-" gorm:"not null;type:text"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"not null;type:timestamptz"`
	VerifiedAt *time.Time `json:"verified_at" gorm:"type:timestamptz"`
	CreatedAt  time.Time  `json:"created_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the EmailVerification model
func (EmailVerification) TableName() string {
	return "email_verifications"
}
