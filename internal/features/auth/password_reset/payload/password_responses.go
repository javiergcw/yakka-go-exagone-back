package payload

// PasswordResetResponse represents a password reset response
type PasswordResetResponse struct {
	Message string `json:"message"`
}

// RequestPasswordResetResponse represents a password reset request response
type RequestPasswordResetResponse struct {
	Message string `json:"message"`
}
