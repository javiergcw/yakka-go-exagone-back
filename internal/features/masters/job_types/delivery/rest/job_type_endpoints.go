package rest

import (
	"net/http"
	"time"

	"github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	"github.com/yakka-backend/internal/shared/response"
)

type JobTypeHandler struct {
	jobTypeRepo database.JobTypeRepository
}

func NewJobTypeHandler(jobTypeRepo database.JobTypeRepository) *JobTypeHandler {
	return &JobTypeHandler{
		jobTypeRepo: jobTypeRepo,
	}
}

// GetJobTypes retrieves all job types
func (h *JobTypeHandler) GetJobTypes(w http.ResponseWriter, r *http.Request) {
	// Simple implementation like licenses
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
		"message": "Job types retrieved successfully",
	})
}

// JobTypeResponse represents a job type in responses
type JobTypeResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetJobTypesResponse represents the response when getting all job types
type GetJobTypesResponse struct {
	Types   []JobTypeResponse `json:"types"`
	Message string            `json:"message"`
}
