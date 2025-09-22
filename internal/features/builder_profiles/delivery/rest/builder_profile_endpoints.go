package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/builder_profiles/payload"
	"github.com/yakka-backend/internal/features/builder_profiles/usecase"
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
	userIDStr, ok := r.Context().Value("user_id").(string)
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
		response.WriteError(w, http.StatusInternalServerError, "Failed to create builder profile")
		return
	}

	// Convert to response
	profileResp := payload.BuilderProfileResponse{
		ID:          profile.ID.String(),
		UserID:      profile.UserID.String(),
		CompanyName: *profile.CompanyName,
		DisplayName: *profile.DisplayName,
		Location:    *profile.Location,
		Bio:         profile.Bio,
		AvatarURL:   profile.AvatarURL,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}

	resp := payload.CreateBuilderProfileResponse{
		Profile: profileResp,
		Message: "Builder profile created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

