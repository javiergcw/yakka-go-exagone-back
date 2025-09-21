package payload

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	SessionID *string `json:"session_id,omitempty"`
}
