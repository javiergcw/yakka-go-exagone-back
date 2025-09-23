package rest

import (
	"context"
	"net/http"

	skillRepo "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	"github.com/yakka-backend/internal/infrastructure/database"
	"github.com/yakka-backend/internal/shared/response"
)

type SkillCompleteHandler struct {
	skillCategoryRepository    skillRepo.SkillCategoryRepository
	skillSubcategoryRepository skillRepo.SkillSubcategoryRepository
}

func NewSkillCompleteHandler() *SkillCompleteHandler {
	return &SkillCompleteHandler{
		skillCategoryRepository:    skillRepo.NewSkillCategoryRepository(database.DB),
		skillSubcategoryRepository: skillRepo.NewSkillSubcategoryRepository(database.DB),
	}
}

// SkillSubcategoryInfo represents subcategory information
type SkillSubcategoryInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SkillCategoryWithSubcategories represents a category with its subcategories
type SkillCategoryWithSubcategories struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	Subcategories []SkillSubcategoryInfo `json:"subcategories"`
}

// GetSkillsComplete returns all skill categories with their subcategories
func (h *SkillCompleteHandler) GetSkillsComplete(w http.ResponseWriter, r *http.Request) {
	// Get all categories
	categories, err := h.skillCategoryRepository.GetAll(context.TODO())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve skill categories")
		return
	}

	// Get all subcategories
	subcategories, err := h.skillSubcategoryRepository.GetAll(context.TODO())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve skill subcategories")
		return
	}

	// Create a map of category ID to subcategories
	subcategoryMap := make(map[string][]SkillSubcategoryInfo)
	for _, subcategory := range subcategories {
		categoryID := subcategory.CategoryID.String()
		subcategoryMap[categoryID] = append(subcategoryMap[categoryID], SkillSubcategoryInfo{
			ID:   subcategory.ID.String(),
			Name: subcategory.Name,
		})
	}

	// Build response with categories and their subcategories
	var result []SkillCategoryWithSubcategories
	for _, category := range categories {
		categoryID := category.ID.String()
		subcategories := subcategoryMap[categoryID]
		if subcategories == nil {
			subcategories = []SkillSubcategoryInfo{} // Empty slice instead of nil
		}

		result = append(result, SkillCategoryWithSubcategories{
			ID:            category.ID.String(),
			Name:          category.Name,
			Description:   category.Description,
			CreatedAt:     category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Subcategories: subcategories,
		})
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
		"message": "Skills with categories and subcategories retrieved successfully",
	})
}
