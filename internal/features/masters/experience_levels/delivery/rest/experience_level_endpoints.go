package rest

import (
	"context"
	"net/http"

	experienceLevelRepo "github.com/yakka-backend/internal/features/masters/experience_levels/entity/database"
	"github.com/yakka-backend/internal/infrastructure/database"
	"github.com/yakka-backend/internal/shared/response"
)

type ExperienceLevelHandler struct {
	experienceLevelRepository experienceLevelRepo.ExperienceLevelRepository
}

func NewExperienceLevelHandler() *ExperienceLevelHandler {
	return &ExperienceLevelHandler{
		experienceLevelRepository: experienceLevelRepo.NewExperienceLevelRepository(database.DB),
	}
}

// GetExperienceLevels returns all experience levels
func (h *ExperienceLevelHandler) GetExperienceLevels(w http.ResponseWriter, r *http.Request) {
	experienceLevels, err := h.experienceLevelRepository.GetAll(context.TODO())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve experience levels")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    experienceLevels,
		"message": "Experience levels retrieved successfully",
	})
}
