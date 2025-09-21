package payload

import "time"

// BuilderProfileResponse represents the response for a builder profile
type BuilderProfileResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	CompanyName string     `json:"company_name"`
	DisplayName string     `json:"display_name"`
	Location    string     `json:"location"`
	Bio         *string    `json:"bio"`
	AvatarURL   *string    `json:"avatar_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateBuilderProfileResponse represents the response when creating a builder profile
type CreateBuilderProfileResponse struct {
	Profile BuilderProfileResponse `json:"profile"`
	Message string                 `json:"message"`
}
