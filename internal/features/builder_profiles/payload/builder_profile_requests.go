package payload

// CreateBuilderProfileRequest represents the request to create/update a builder profile
type CreateBuilderProfileRequest struct {
	CompanyName string  `json:"company_name" validate:"required,min=2,max=255"`
	DisplayName string  `json:"display_name" validate:"required,min=2,max=255"`
	Location    string  `json:"location" validate:"required,min=2,max=255"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,min=10,max=32"`
}
