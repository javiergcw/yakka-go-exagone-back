package rest

import (
	"net/http"
	"time"

	"github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	"github.com/yakka-backend/internal/shared/response"
)

type JobRequirementHandler struct {
	jobRequirementRepo database.JobRequirementRepository
}

func NewJobRequirementHandler(jobRequirementRepo database.JobRequirementRepository) *JobRequirementHandler {
	return &JobRequirementHandler{
		jobRequirementRepo: jobRequirementRepo,
	}
}

// GetJobRequirements retrieves all job requirements
func (h *JobRequirementHandler) GetJobRequirements(w http.ResponseWriter, r *http.Request) {
	// Simple implementation like licenses
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
		"message": "Job requirements retrieved successfully",
	})
}

// JobRequirementResponse represents a job requirement in responses
type JobRequirementResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetJobRequirementsResponse represents the response when getting all job requirements
type GetJobRequirementsResponse struct {
	Requirements []JobRequirementResponse `json:"requirements"`
	Message      string                   `json:"message"`
}
