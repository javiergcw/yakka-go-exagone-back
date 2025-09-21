package models

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
	UserStatusBanned   UserStatus = "banned"
)

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleUser    UserRole = "user"
	UserRoleBuilder UserRole = "builder"
	UserRoleLabour  UserRole = "labour"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string     `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Phone         *string    `json:"phone" gorm:"size:32"`
	PasswordHash  string     `json:"-" gorm:"not null;type:text"`
	Status        UserStatus `json:"status" gorm:"not null;type:user_status"`
	LastLoginAt   *time.Time `json:"last_login_at" gorm:"type:timestamptz"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null;type:timestamptz"`
	Role          UserRole   `json:"role" gorm:"not null;type:user_role"`
	RoleChangedAt *time.Time `json:"role_changed_at" gorm:"type:timestamptz"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}
