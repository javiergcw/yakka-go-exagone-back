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
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}
