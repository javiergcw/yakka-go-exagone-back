package payload

// RequestEmailVerificationRequest represents an email verification request
type RequestEmailVerificationRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

// VerifyEmailRequest represents an email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}
