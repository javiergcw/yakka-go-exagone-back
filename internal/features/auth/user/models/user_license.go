package models

import (
	"time"

	"github.com/google/uuid"
	license_models "github.com/yakka-backend/internal/features/masters/licenses/models"
)

// UserLicense represents a license associated with a user
type UserLicense struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_user_license"`
	LicenseID uuid.UUID  `json:"license_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_user_license"`
	PhotoURL  *string    `json:"photo_url" gorm:"type:text"`
	IssuedAt  *time.Time `json:"issued_at" gorm:"type:timestamptz"`
	ExpiresAt *time.Time `json:"expires_at" gorm:"type:timestamptz"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;type:timestamptz"`

	// Relationships
	User    User                   `json:"user" gorm:"foreignKey:UserID;references:ID"`
	License license_models.License `json:"license" gorm:"foreignKey:LicenseID;references:ID"`
}

// TableName returns the table name for the UserLicense model
func (UserLicense) TableName() string {
	return "user_licenses"
}
