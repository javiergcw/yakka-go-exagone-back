package rest

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	skillRepo "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	"github.com/yakka-backend/internal/infrastructure/database"
	"github.com/yakka-backend/internal/shared/response"
)

type SkillSubcategoryHandler struct {
	skillSubcategoryRepository skillRepo.SkillSubcategoryRepository
}

func NewSkillSubcategoryHandler() *SkillSubcategoryHandler {
	return &SkillSubcategoryHandler{
		skillSubcategoryRepository: skillRepo.NewSkillSubcategoryRepository(database.DB),
	}
}

// GetSkillSubcategories returns all skill subcategories or filtered by category
func (h *SkillSubcategoryHandler) GetSkillSubcategories(w http.ResponseWriter, r *http.Request) {
	// Check if category_id query parameter is provided
	categoryID := r.URL.Query().Get("category_id")

	var subcategories interface{}
	var err error

	if categoryID != "" {
		// Parse categoryID to UUID
		categoryUUID, parseErr := uuid.Parse(categoryID)
		if parseErr != nil {
			response.WriteError(w, http.StatusBadRequest, "Invalid category ID format")
			return
		}
		// Get subcategories for specific category
		subcategories, err = h.skillSubcategoryRepository.GetByCategoryID(context.TODO(), categoryUUID)
	} else {
		// Get all subcategories
		subcategories, err = h.skillSubcategoryRepository.GetAll(context.TODO())
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve skill subcategories")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    subcategories,
		"message": "Skill subcategories retrieved successfully",
	})
}

// GetSkillSubcategoriesByCategory returns subcategories for a specific category
func (h *SkillSubcategoryHandler) GetSkillSubcategoriesByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["categoryId"]

	if categoryID == "" {
		response.WriteError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	// Parse categoryID to UUID
	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	subcategories, err := h.skillSubcategoryRepository.GetByCategoryID(context.TODO(), categoryUUID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve skill subcategories")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    subcategories,
		"message": "Skill subcategories retrieved successfully",
	})
}
