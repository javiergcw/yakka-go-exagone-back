package payload

// CreateLabourProfileRequest represents the request to create/update a labour profile
type CreateLabourProfileRequest struct {
	FirstName string                      `json:"first_name" validate:"required,min=2,max=120"`
	LastName  string                      `json:"last_name" validate:"required,min=2,max=120"`
	Location  string                      `json:"location" validate:"required,min=2,max=255"`
	Bio       *string                     `json:"bio,omitempty"`
	AvatarURL *string                     `json:"avatar_url,omitempty"`
	Phone     *string                     `json:"phone,omitempty" validate:"omitempty,min=10,max=32"`
	Skills    []LabourProfileSkillRequest `json:"skills,omitempty"`
	Licenses  []UserLicenseRequest        `json:"licenses,omitempty"`
}

// LabourProfileSkillRequest represents a skill to be added to a labour profile
type LabourProfileSkillRequest struct {
	CategoryID        string  `json:"category_id" validate:"required,uuid"`
	SubcategoryID     string  `json:"subcategory_id" validate:"required,uuid"`
	ExperienceLevelID string  `json:"experience_level_id" validate:"required,uuid"`
	YearsExperience   float64 `json:"years_experience" validate:"required,min=0,max=99.9"`
	IsPrimary         bool    `json:"is_primary"`
}

// UserLicenseRequest represents a license to be added to a user profile
type UserLicenseRequest struct {
	LicenseID string  `json:"license_id" validate:"required,uuid"`
	PhotoURL  *string `json:"photo_url,omitempty"`
	IssuedAt  *string `json:"issued_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ExpiresAt *string `json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}
