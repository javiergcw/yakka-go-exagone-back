package rest

import (
	"net/http"
	"time"

	"github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	"github.com/yakka-backend/internal/features/masters/job_types/models"
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
	// Check if only active types are requested
	activeOnly := r.URL.Query().Get("active_only")

	var types []*models.JobType
	var err error

	if activeOnly == "true" {
		types, err = h.jobTypeRepo.GetActive(r.Context())
	} else {
		types, err = h.jobTypeRepo.GetAll(r.Context())
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get job types")
		return
	}

	// Convert to response
	var typesResp []JobTypeResponse
	for _, jobType := range types {
		typesResp = append(typesResp, JobTypeResponse{
			ID:          jobType.ID.String(),
			Name:        jobType.Name,
			Description: jobType.Description,
			IsActive:    jobType.IsActive,
			CreatedAt:   jobType.CreatedAt,
			UpdatedAt:   jobType.UpdatedAt,
		})
	}

	resp := GetJobTypesResponse{
		Types:   typesResp,
		Message: "Job types retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
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
