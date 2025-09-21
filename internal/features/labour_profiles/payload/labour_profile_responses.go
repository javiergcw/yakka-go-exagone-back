package payload

import "time"

// LabourProfileResponse represents the response for a labour profile
type LabourProfileResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Location  string     `json:"location"`
	Bio       *string    `json:"bio"`
	AvatarURL *string    `json:"avatar_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CreateLabourProfileResponse represents the response when creating a labour profile
type CreateLabourProfileResponse struct {
	Profile LabourProfileResponse `json:"profile"`
	Message string                `json:"message"`
}
