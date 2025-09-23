package payload

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email      string  `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Phone      *string `json:"phone,omitempty"`
	AutoVerify bool   `json:"auto_verify,omitempty"` // Si true, verifica la cuenta inmediatamente
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Phone     *string `json:"phone,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Address   *string `json:"address,omitempty"`
	Photo     *string `json:"photo,omitempty"`
}
