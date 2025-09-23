package payload

import "time"

// UserResponse represents a user response
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         *string    `json:"phone"`
	FirstName     *string    `json:"first_name"`
	LastName      *string    `json:"last_name"`
	Address       *string    `json:"address"`
	Photo         *string    `json:"photo"`
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
	FirstName     *string    `json:"first_name"`
	LastName      *string    `json:"last_name"`
	Address       *string    `json:"address"`
	Photo         *string    `json:"photo"`
	Status        string     `json:"status"`
	Role          string     `json:"role"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	RoleChangedAt *time.Time `json:"role_changed_at"`
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	User         RegisterUserResponse `json:"user"`
	Message      string               `json:"message"`
	AutoVerified bool                 `json:"auto_verified,omitempty"` // Indica si fue verificado automáticamente
	EmailSent    bool                 `json:"email_sent,omitempty"`    // Indica si se envió email de verificación
}

// CompleteProfileResponse represents a complete user profile response with all profiles
type CompleteProfileResponse struct {
	User              UserResponse        `json:"user"`
	BuilderProfile    *BuilderProfileInfo `json:"builder_profile,omitempty"`
	LabourProfile     *LabourProfileInfo  `json:"labour_profile,omitempty"`
	CurrentRole       string              `json:"current_role"`
	HasBuilderProfile bool                `json:"has_builder_profile"`
	HasLabourProfile  bool                `json:"has_labour_profile"`
}

// BuilderProfileInfo represents builder profile information
type BuilderProfileInfo struct {
	ID          string  `json:"id"`
	CompanyName *string `json:"company_name"`
	DisplayName *string `json:"display_name"`
	Location    *string `json:"location"`
	Bio         *string `json:"bio"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// LabourProfileInfo represents labour profile information
type LabourProfileInfo struct {
	ID        string  `json:"id"`
	Location  *string `json:"location"`
	Bio       *string `json:"bio"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
