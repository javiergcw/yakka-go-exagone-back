package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/http/middleware"
	"github.com/yakka-backend/internal/features/auth/email_verification/usecase"
	"github.com/yakka-backend/internal/features/auth/user/payload"
	auth_user_usecase "github.com/yakka-backend/internal/features/auth/user/usecase"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authUsecase              auth_user_usecase.AuthUsecase
	emailVerificationUsecase usecase.EmailVerificationUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase auth_user_usecase.AuthUsecase, emailVerificationUsecase usecase.EmailVerificationUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase:              authUsecase,
		emailVerificationUsecase: emailVerificationUsecase,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req payload.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Register user
	user, err := h.authUsecase.Register(r.Context(), req.Email, req.Password, req.Phone)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Conflict":
			response.WriteError(w, http.StatusConflict, "User with this email already exists")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid registration data")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	// Send email verification automatically
	token, err := h.emailVerificationUsecase.RequestEmailVerification(r.Context(), user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to send email verification to %s: %v", user.Email, err)
		// Continue anyway - user is created, just email verification failed
	} else {
		log.Printf("📧 Email verification sent to %s with token: %s", user.Email, token)
	}

	// Convert to response
	userResp := payload.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		Status:        string(user.Status),
		Role:          string(user.Role),
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		RoleChangedAt: user.RoleChangedAt,
	}

	resp := payload.RegisterResponse{
		User:    userResp,
		Message: "User registered successfully. Please check your email for verification.",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req payload.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Login user
	user, err := h.authUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Forbidden":
			response.WriteError(w, http.StatusForbidden, "Email verification required. Please check your email and click the verification link before logging in.")
		case "Unauthorized":
			response.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		case "Not found":
			response.WriteError(w, http.StatusUnauthorized, "User not found")
		default:
			response.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		}
		return
	}

	// Convert to response
	userResp := payload.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		Status:        string(user.Status),
		Role:          string(user.Role),
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		RoleChangedAt: user.RoleChangedAt,
	}

	// Generate JWT access token
	accessToken, err := middleware.GenerateJWTToken(user.ID.String())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// TODO: Generate refresh token (for now using placeholder)
	refreshToken := "refresh_token_here"
	expiresIn := int64(3600) // 1 hour

	resp := payload.LoginResponse{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get user
	user, err := h.authUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User not found")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve user")
		}
		return
	}

	// Convert to response
	userResp := payload.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		Status:        string(user.Status),
		Role:          string(user.Role),
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		RoleChangedAt: user.RoleChangedAt,
	}

	response.WriteJSON(w, http.StatusOK, userResp)
}

// UpdateProfile handles updating user profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req payload.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get current user
	user, err := h.authUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User not found")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve user")
		}
		return
	}

	// Update fields
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	// Update user
	err = h.authUsecase.UpdateUser(r.Context(), user)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User not found")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid update data")
		default:
			response.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	// Convert to response
	userResp := payload.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		Status:        string(user.Status),
		Role:          string(user.Role),
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		RoleChangedAt: user.RoleChangedAt,
	}

	response.WriteJSON(w, http.StatusOK, userResp)
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req payload.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Change password
	err = h.authUsecase.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		// Handle specific error types
		switch err.Error() {
		case "Not found":
			response.WriteError(w, http.StatusNotFound, "User not found")
		case "Unauthorized":
			response.WriteError(w, http.StatusUnauthorized, "Current password is incorrect")
		case "Bad request":
			response.WriteError(w, http.StatusBadRequest, "Invalid password data")
		default:
			response.WriteError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}
