package payload

// CreateJobsiteRequest represents the request to create a jobsite
type CreateJobsiteRequest struct {
	BuilderID   string  `json:"builder_id" validate:"required,uuid"`
	Address     string  `json:"address" validate:"required,min=10,max=500"`
	City        *string `json:"city,omitempty" validate:"omitempty,min=2,max=120"`
	Suburb      *string `json:"suburb,omitempty" validate:"omitempty,min=2,max=120"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,min=10,max=32"`
}

// UpdateJobsiteRequest represents the request to update a jobsite
type UpdateJobsiteRequest struct {
	Address     *string  `json:"address,omitempty" validate:"omitempty,min=10,max=500"`
	City        *string  `json:"city,omitempty" validate:"omitempty,min=2,max=120"`
	Suburb      *string  `json:"suburb,omitempty" validate:"omitempty,min=2,max=120"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Latitude    *float64 `json:"latitude,omitempty" validate:"omitempty,min=-90,max=90"`
	Longitude   *float64 `json:"longitude,omitempty" validate:"omitempty,min=-180,max=180"`
	Phone       *string  `json:"phone,omitempty" validate:"omitempty,min=10,max=32"`
}
