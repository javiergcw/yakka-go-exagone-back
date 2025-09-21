package payload

// CreateLabourProfileRequest represents the request to create/update a labour profile
type CreateLabourProfileRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=2,max=120"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=120"`
	Location  string  `json:"location" validate:"required,min=2,max=255"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,min=10,max=32"`
}
