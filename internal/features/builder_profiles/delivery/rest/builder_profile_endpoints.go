package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/builder_profiles/payload"
	"github.com/yakka-backend/internal/features/builder_profiles/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type BuilderProfileHandler struct {
	builderProfileUsecase usecase.BuilderProfileUsecase
}

func NewBuilderProfileHandler(builderProfileUsecase usecase.BuilderProfileUsecase) *BuilderProfileHandler {
	return &BuilderProfileHandler{
		builderProfileUsecase: builderProfileUsecase,
	}
}

// CreateBuilderProfile creates or updates a builder profile
func (h *BuilderProfileHandler) CreateBuilderProfile(w http.ResponseWriter, r *http.Request) {
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

	var req payload.CreateBuilderProfileRequest
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
	profile, err := h.builderProfileUsecase.CreateProfile(r.Context(), userID, req)
	if err != nil {
		// Check for specific error types
		if err.Error() == "builder profile already exists for this user" {
			response.WriteError(w, http.StatusConflict, "Builder profile already exists for this user")
			return
		}

		// Check for validation errors (UUID not found)
		if err.Error() == "license not found" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check for invalid UUID format
		if err.Error() == "invalid license_id" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check for company not found
		if err.Error() == "company not found" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check for invalid company_id format
		if err.Error() == "invalid company_id" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generic error for other cases
		response.WriteError(w, http.StatusInternalServerError, "Failed to create builder profile")
		return
	}

	// Convert to response
	profileResp := payload.BuilderProfileResponse{
		ID:          profile.ID.String(),
		UserID:      profile.UserID.String(),
		DisplayName: getStringValue(profile.DisplayName),
		Location:    getStringValue(profile.Location),
		Bio:         profile.Bio,
		AvatarURL:   req.AvatarURL, // Usar datos del request ya que están en el usuario
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}

	// Add company information if available
	if profile.CompanyID != nil {
		companyIDStr := profile.CompanyID.String()
		profileResp.CompanyID = &companyIDStr

		// If company data is loaded, include it in response
		if profile.Company != nil {
			profileResp.Company = &payload.CompanyResponse{
				ID:          profile.Company.ID.String(),
				Name:        profile.Company.Name,
				Description: profile.Company.Description,
				Website:     profile.Company.Website,
			}
		}
	}

	resp := payload.CreateBuilderProfileResponse{
		Profile: profileResp,
		Message: "Builder profile created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
