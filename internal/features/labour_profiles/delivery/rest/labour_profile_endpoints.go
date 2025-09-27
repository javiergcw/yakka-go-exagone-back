package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/labour_profiles/payload"
	"github.com/yakka-backend/internal/features/labour_profiles/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type LabourProfileHandler struct {
	labourProfileUsecase usecase.LabourProfileUsecase
}

func NewLabourProfileHandler(labourProfileUsecase usecase.LabourProfileUsecase) *LabourProfileHandler {
	return &LabourProfileHandler{
		labourProfileUsecase: labourProfileUsecase,
	}
}

// CreateLabourProfile creates or updates a labour profile
func (h *LabourProfileHandler) CreateLabourProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req payload.CreateLabourProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create profile
	profile, err := h.labourProfileUsecase.CreateProfile(r.Context(), userID, req)
	if err != nil {
		// Check for specific error types
		if err.Error() == "labour profile already exists for this user" {
			response.WriteError(w, http.StatusConflict, "Labour profile already exists for this user")
			return
		}

		// Check for validation errors (UUID not found)
		if err.Error() == "skill category not found" ||
			err.Error() == "skill subcategory not found" ||
			err.Error() == "experience level not found" ||
			err.Error() == "license not found" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check for invalid UUID format
		if err.Error() == "invalid category_id" ||
			err.Error() == "invalid subcategory_id" ||
			err.Error() == "invalid experience_level_id" ||
			err.Error() == "invalid license_id" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generic error for other cases - include actual error for debugging
		response.WriteError(w, http.StatusInternalServerError, "Failed to create labour profile: "+err.Error())
		return
	}

	// Convert to response
	profileResp := payload.LabourProfileResponse{
		ID:        profile.ID.String(),
		UserID:    profile.UserID.String(),
		FirstName: req.FirstName, // Usar datos del request ya que están en el usuario
		LastName:  req.LastName,  // Usar datos del request ya que están en el usuario
		Location:  *profile.Location,
		Bio:       profile.Bio,
		AvatarURL: req.AvatarURL, // Usar datos del request ya que están en el usuario
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	resp := payload.CreateLabourProfileResponse{
		Profile: profileResp,
		Message: "Labour profile created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}
