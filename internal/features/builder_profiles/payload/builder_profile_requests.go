package payload

// CreateBuilderProfileRequest represents the request to create/update a builder profile
type CreateBuilderProfileRequest struct {
	CompanyName string               `json:"company_name" validate:"required,min=2,max=255"`
	DisplayName string               `json:"display_name" validate:"required,min=2,max=255"`
	Location    string               `json:"location" validate:"required,min=2,max=255"`
	Bio         *string              `json:"bio,omitempty"`
	AvatarURL   *string              `json:"avatar_url,omitempty"`
	Phone       *string              `json:"phone,omitempty" validate:"omitempty,min=10,max=32"`
	Licenses    []UserLicenseRequest `json:"licenses,omitempty"`
}

// UserLicenseRequest represents a license to be added to a user profile
type UserLicenseRequest struct {
	LicenseID string  `json:"license_id" validate:"required,uuid"`
	PhotoURL  *string `json:"photo_url,omitempty"`
	IssuedAt  *string `json:"issued_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ExpiresAt *string `json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}
