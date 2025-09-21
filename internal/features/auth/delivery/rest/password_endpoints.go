package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yakka-backend/internal/features/auth/password_reset/payload"
	"github.com/yakka-backend/internal/features/auth/password_reset/usecase"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

// PasswordHandler handles password reset endpoints
type PasswordHandler struct {
	passwordResetUsecase usecase.PasswordResetUsecase
}

// NewPasswordHandler creates a new password handler
func NewPasswordHandler(passwordResetUsecase usecase.PasswordResetUsecase) *PasswordHandler {
	return &PasswordHandler{
		passwordResetUsecase: passwordResetUsecase,
	}
}

// RequestPasswordReset handles password reset request
func (h *PasswordHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req payload.RequestPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Request password reset
	token, err := h.passwordResetUsecase.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User with this email not found")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid email address")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to request password reset")
		}
		return
	}

	// Send email with reset token
	// TODO: Implement actual email sending service
	// For now, assume all emails are sent successfully
	log.Printf("📧 Password reset email sent to %s with token: %s", req.Email, token)

	resp := payload.RequestPasswordResetResponse{
		Message: "Password reset email sent",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// ResetPassword handles password reset
func (h *PasswordHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req payload.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Reset password
	err := h.passwordResetUsecase.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Unauthorized":
			response.WriteError(w, http.StatusUnauthorized, "Invalid or expired reset token")
		case "Conflict":
			response.WriteError(w, http.StatusConflict, "Reset token already used")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid password or token data")
		default:
			response.WriteError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	resp := payload.PasswordResetResponse{
		Message: "Password reset successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
