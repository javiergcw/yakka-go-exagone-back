package payload

import "time"

// UserResponse represents a user response
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         *string    `json:"phone"`
	Status        string     `json:"status"`
	Role          string     `json:"role"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	RoleChangedAt *time.Time `json:"role_changed_at"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	Profiles     ProfileInfo  `json:"profiles"`
}

// ProfileInfo represents information about user profiles
type ProfileInfo struct {
	HasBuilderProfile bool `json:"has_builder_profile"`
	HasLabourProfile  bool `json:"has_labour_profile"`
	HasAnyProfile     bool `json:"has_any_profile"`
}

// RegisterUserResponse represents a user response for registration (without phone)
type RegisterUserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Status        string     `json:"status"`
	Role          string     `json:"role"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	RoleChangedAt *time.Time `json:"role_changed_at"`
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	User           RegisterUserResponse `json:"user"`
	Message        string                `json:"message"`
	AutoVerified   bool                  `json:"auto_verified,omitempty"`   // Indica si fue verificado automáticamente
	EmailSent      bool                  `json:"email_sent,omitempty"`      // Indica si se envió email de verificación
}
