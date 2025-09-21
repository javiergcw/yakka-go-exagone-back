package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/email_verification/payload"
	"github.com/yakka-backend/internal/features/auth/email_verification/usecase"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

// EmailHandler handles email verification endpoints
type EmailHandler struct {
	emailVerificationUsecase usecase.EmailVerificationUsecase
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(emailVerificationUsecase usecase.EmailVerificationUsecase) *EmailHandler {
	return &EmailHandler{
		emailVerificationUsecase: emailVerificationUsecase,
	}
}

// RequestEmailVerification handles email verification request
func (h *EmailHandler) RequestEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req payload.RequestEmailVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Request email verification
	token, err := h.emailVerificationUsecase.RequestEmailVerification(r.Context(), userID)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User not found")
		case "Conflict":
			response.WriteError(w, http.StatusConflict, "Email verification already requested")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid user data")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to request email verification")
		}
		return
	}

	// Send email with verification token
	// TODO: Implement actual email sending service
	// For now, assume all emails are verified automatically
	log.Printf("📧 Email verification sent to user %s with token: %s", req.UserID, token)

	resp := payload.RequestEmailVerificationResponse{
		Message: "Verification email sent",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// VerifyEmail handles email verification
func (h *EmailHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req payload.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Verify email
	err := h.emailVerificationUsecase.VerifyEmail(r.Context(), req.Token)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Unauthorized":
			response.WriteError(w, http.StatusUnauthorized, "Invalid or expired verification token")
		case "Conflict":
			response.WriteError(w, http.StatusConflict, "Email already verified")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid verification token")
		default:
			response.WriteError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	resp := payload.EmailVerificationResponse{
		Message: "Email verified successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
