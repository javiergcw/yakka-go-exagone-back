package payload

// EmailVerificationResponse represents an email verification response
type EmailVerificationResponse struct {
	Message string `json:"message"`
}

// RequestEmailVerificationResponse represents an email verification request response
type RequestEmailVerificationResponse struct {
	Message string `json:"message"`
}
