package payload

import "time"

// BuilderProfileResponse represents the response for a builder profile
type BuilderProfileResponse struct {
	ID          string           `json:"id"`
	UserID      string           `json:"user_id"`
	CompanyID   *string          `json:"company_id,omitempty"`
	Company     *CompanyResponse `json:"company,omitempty"`
	DisplayName string           `json:"display_name"`
	Location    string           `json:"location"`
	Bio         *string          `json:"bio"`
	AvatarURL   *string          `json:"avatar_url"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// CompanyResponse represents the response for a company
type CompanyResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Website     *string   `json:"website"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateBuilderProfileResponse represents the response when creating a builder profile
type CreateBuilderProfileResponse struct {
	Profile BuilderProfileResponse `json:"profile"`
	Message string                 `json:"message"`
}

// CreateCompanyResponse represents the response when creating a company
type CreateCompanyResponse struct {
	Company CompanyResponse `json:"company"`
	Message string          `json:"message"`
}

// GetCompaniesResponse represents the response when getting all companies
type GetCompaniesResponse struct {
	Companies []CompanyResponse `json:"companies"`
	Total     int               `json:"total"`
	Message   string            `json:"message"`
}

// AssignCompanyResponse represents the response when assigning a company to a builder
type AssignCompanyResponse struct {
	BuilderProfile BuilderProfileResponse `json:"builder_profile"`
	Message        string                 `json:"message"`
}
