package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/email_verification/usecase"
	"github.com/yakka-backend/internal/features/auth/user/models"
	"github.com/yakka-backend/internal/features/auth/user/payload"
	auth_user_usecase "github.com/yakka-backend/internal/features/auth/user/usecase"
	builder_usecase "github.com/yakka-backend/internal/features/builder_profiles/usecase"
	labour_usecase "github.com/yakka-backend/internal/features/labour_profiles/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
	"gorm.io/gorm"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authUsecase              auth_user_usecase.AuthUsecase
	emailVerificationUsecase usecase.EmailVerificationUsecase
	builderProfileUsecase    builder_usecase.BuilderProfileUsecase
	labourProfileUsecase     labour_usecase.LabourProfileUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	authUsecase auth_user_usecase.AuthUsecase,
	emailVerificationUsecase usecase.EmailVerificationUsecase,
	builderProfileUsecase builder_usecase.BuilderProfileUsecase,
	labourProfileUsecase labour_usecase.LabourProfileUsecase,
) *AuthHandler {
	return &AuthHandler{
		authUsecase:              authUsecase,
		emailVerificationUsecase: emailVerificationUsecase,
		builderProfileUsecase:    builderProfileUsecase,
		labourProfileUsecase:     labourProfileUsecase,
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

	// Register user based on auto_verify parameter
	var user *models.User
	var err error
	var isAutoVerified bool

	if req.AutoVerify {
		// Register with auto verification
		user, err = h.authUsecase.RegisterWithAutoVerify(r.Context(), req.Email, req.Password, req.Phone)
		isAutoVerified = true
	} else {
		// Register with pending status (requires email verification)
		user, err = h.authUsecase.Register(r.Context(), req.Email, req.Password, req.Phone)
		isAutoVerified = false
	}

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

	// Send email verification only if not auto-verified
	var emailSent bool
	if !isAutoVerified {
		token, err := h.emailVerificationUsecase.RequestEmailVerification(r.Context(), user.ID)
		if err != nil {
			log.Printf("⚠️ Failed to send email verification to %s: %v", user.Email, err)
			emailSent = false
		} else {
			log.Printf("📧 Email verification sent to %s with token: %s", user.Email, token)
			emailSent = true
		}
	}

	// Convert to response (without phone for registration)
	userResp := payload.RegisterUserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Address:       user.Address,
		Photo:         user.Photo,
		Status:        string(user.Status),
		Role:          string(user.Role),
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		RoleChangedAt: user.RoleChangedAt,
	}

	// Set appropriate message based on verification status
	var message string
	if isAutoVerified {
		message = "User registered and verified successfully. You can now log in."
	} else {
		message = "User registered successfully. Please check your email for verification."
	}

	resp := payload.RegisterResponse{
		User:         userResp,
		Message:      message,
		AutoVerified: isAutoVerified,
		EmailSent:    emailSent,
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
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Address:       user.Address,
		Photo:         user.Photo,
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
		log.Printf("❌ Failed to generate JWT token: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Printf("✅ JWT token generated successfully for user: %s", user.ID.String())

	// Get user profile information
	profileInfo, err := h.authUsecase.GetUserProfileInfo(r.Context(), user.ID)
	if err != nil {
		// If we can't get profile info, set default values
		profileInfo = payload.ProfileInfo{
			HasBuilderProfile: false,
			HasLabourProfile:  false,
			HasAnyProfile:     false,
		}
	}

	// TODO: Generate refresh token (for now using placeholder)
	refreshToken := "refresh_token_here"
	expiresIn := int64(3600) // 1 hour

	resp := payload.LoginResponse{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		Profiles:     profileInfo,
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetProfile handles getting complete user profile with all profiles
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr := r.Context().Value(middleware.UserIDKey).(string)
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

	// Convert user to response
	userResp := payload.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Address:       user.Address,
		Photo:         user.Photo,
		Status:        string(user.Status),
		Role:          string(user.Role),
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		RoleChangedAt: user.RoleChangedAt,
	}

	// Initialize response
	completeResp := payload.CompleteProfileResponse{
		User:              userResp,
		CurrentRole:       string(user.Role),
		HasBuilderProfile: false,
		HasLabourProfile:  false,
		BuilderProfile:    nil,
		LabourProfile:     nil,
	}

	// Try to get builder profile
	builderProfile, err := h.builderProfileUsecase.GetProfileByUserID(r.Context(), userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error getting builder profile: %v", err)
		// Set empty builder profile on error
		completeResp.BuilderProfile = &payload.BuilderProfileInfo{}
	} else if err == nil && builderProfile != nil {
		completeResp.HasBuilderProfile = true
		completeResp.BuilderProfile = &payload.BuilderProfileInfo{
			ID:          builderProfile.ID.String(),
			CompanyName: builderProfile.CompanyName,
			DisplayName: builderProfile.DisplayName,
			Location:    builderProfile.Location,
			Bio:         builderProfile.Bio,
			CreatedAt:   builderProfile.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   builderProfile.UpdatedAt.Format(time.RFC3339),
		}
	} else {
		// No builder profile found, set empty
		completeResp.BuilderProfile = &payload.BuilderProfileInfo{}
	}

	// Try to get labour profile
	labourProfile, err := h.labourProfileUsecase.GetProfileByUserID(r.Context(), userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error getting labour profile: %v", err)
		// Set empty labour profile on error
		completeResp.LabourProfile = &payload.LabourProfileInfo{}
	} else if err == nil && labourProfile != nil {
		completeResp.HasLabourProfile = true
		completeResp.LabourProfile = &payload.LabourProfileInfo{
			ID:        labourProfile.ID.String(),
			Location:  labourProfile.Location,
			Bio:       labourProfile.Bio,
			CreatedAt: labourProfile.CreatedAt.Format(time.RFC3339),
			UpdatedAt: labourProfile.UpdatedAt.Format(time.RFC3339),
		}
	} else {
		// No labour profile found, set empty
		completeResp.LabourProfile = &payload.LabourProfileInfo{}
	}

	response.WriteJSON(w, http.StatusOK, completeResp)
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
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.Photo != nil {
		user.Photo = req.Photo
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
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Address:       user.Address,
		Photo:         user.Photo,
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
