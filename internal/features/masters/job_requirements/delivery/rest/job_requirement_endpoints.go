package rest

import (
	"net/http"
	"time"

	"github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	"github.com/yakka-backend/internal/features/masters/job_requirements/models"
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
	// Check if only active requirements are requested
	activeOnly := r.URL.Query().Get("active_only")

	var requirements []*models.JobRequirement
	var err error

	if activeOnly == "true" {
		requirements, err = h.jobRequirementRepo.GetActive(r.Context())
	} else {
		requirements, err = h.jobRequirementRepo.GetAll(r.Context())
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get job requirements")
		return
	}

	// Convert to response
	var requirementsResp []JobRequirementResponse
	for _, requirement := range requirements {
		requirementsResp = append(requirementsResp, JobRequirementResponse{
			ID:          requirement.ID.String(),
			Name:        requirement.Name,
			Description: requirement.Description,
			IsActive:    requirement.IsActive,
			CreatedAt:   requirement.CreatedAt,
			UpdatedAt:   requirement.UpdatedAt,
		})
	}

	resp := GetJobRequirementsResponse{
		Requirements: requirementsResp,
		Message:      "Job requirements retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
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
