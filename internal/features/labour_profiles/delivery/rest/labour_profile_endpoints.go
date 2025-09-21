package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/labour_profiles/payload"
	"github.com/yakka-backend/internal/features/labour_profiles/usecase"
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
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
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
		response.WriteError(w, http.StatusInternalServerError, "Failed to create labour profile")
		return
	}

	// Convert to response
	profileResp := payload.LabourProfileResponse{
		ID:        profile.ID.String(),
		UserID:    profile.UserID.String(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Location:  profile.Location,
		Bio:       profile.Bio,
		AvatarURL: profile.AvatarURL,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	resp := payload.CreateLabourProfileResponse{
		Profile: profileResp,
		Message: "Labour profile created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

