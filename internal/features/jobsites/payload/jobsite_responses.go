package payload

import (
	"time"

	"github.com/google/uuid"
)

// JobsiteResponse represents the response for a jobsite
type JobsiteResponse struct {
	ID          uuid.UUID `json:"id"`
	BuilderID   uuid.UUID `json:"builder_id"`
	Address     string    `json:"address"`
	City        *string   `json:"city,omitempty"`
	Suburb      *string   `json:"suburb,omitempty"`
	Description *string   `json:"description,omitempty"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Phone       *string   `json:"phone,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JobsiteListResponse represents the response for a list of jobsites
type JobsiteListResponse struct {
	Jobsites []JobsiteResponse `json:"jobsites"`
	Total    int               `json:"total"`
}

// CreateJobsiteResponse represents the response after creating a jobsite
type CreateJobsiteResponse struct {
	Jobsite JobsiteResponse `json:"jobsite"`
	Message string          `json:"message"`
}

// UpdateJobsiteResponse represents the response after updating a jobsite
type UpdateJobsiteResponse struct {
	Jobsite JobsiteResponse `json:"jobsite"`
	Message string          `json:"message"`
}

// DeleteJobsiteResponse represents the response after deleting a jobsite
type DeleteJobsiteResponse struct {
	Message string `json:"message"`
}
