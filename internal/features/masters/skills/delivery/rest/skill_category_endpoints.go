package rest

import (
	"context"
	"net/http"

	skillRepo "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	"github.com/yakka-backend/internal/infrastructure/database"
	"github.com/yakka-backend/internal/shared/response"
)

type SkillCategoryHandler struct {
	skillCategoryRepository skillRepo.SkillCategoryRepository
}

func NewSkillCategoryHandler() *SkillCategoryHandler {
	return &SkillCategoryHandler{
		skillCategoryRepository: skillRepo.NewSkillCategoryRepository(database.DB),
	}
}

// GetSkillCategories returns all skill categories
func (h *SkillCategoryHandler) GetSkillCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.skillCategoryRepository.GetAll(context.TODO())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve skill categories")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    categories,
		"message": "Skill categories retrieved successfully",
	})
}
