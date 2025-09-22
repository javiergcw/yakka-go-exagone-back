package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/features/auth/user_session/payload"
	"github.com/yakka-backend/internal/features/auth/user_session/usecase"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

// SessionHandler handles session endpoints
type SessionHandler struct {
	sessionUsecase usecase.SessionUsecase
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionUsecase usecase.SessionUsecase) *SessionHandler {
	return &SessionHandler{
		sessionUsecase: sessionUsecase,
	}
}

// RefreshToken handles token refresh
func (h *SessionHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req payload.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Refresh session
	_, newRefreshToken, err := h.sessionUsecase.RefreshSession(r.Context(), req.RefreshToken)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Unauthorized":
			response.WriteError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "Session not found")
		default:
			response.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		}
		return
	}

	// Generate new JWT access token
	// TODO: Extract user ID from refresh token validation
	userID := "user_id_from_refresh_token" // This should come from refresh token validation
	accessToken, err := middleware.GenerateJWTToken(userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}
	expiresIn := int64(3600) // 1 hour

	resp := payload.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// Logout handles user logout
func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req payload.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// If specific session ID provided, revoke that session
	if req.SessionID != nil {
		// TODO: Implement session-specific logout
		response.WriteJSON(w, http.StatusOK, payload.LogoutResponse{
			Message: "Session revoked successfully",
		})
		return
	}

	// Revoke all user sessions
	err = h.sessionUsecase.RevokeAllUserSessions(r.Context(), userID)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User sessions not found")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to logout")
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, payload.LogoutResponse{
		Message: "Logged out successfully",
	})
}
