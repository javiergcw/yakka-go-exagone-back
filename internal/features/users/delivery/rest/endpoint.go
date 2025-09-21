package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yakka-backend/internal/features/users"
	"github.com/yakka-backend/internal/features/users/models"
	"github.com/yakka-backend/internal/features/users/payload"
	"github.com/yakka-backend/internal/shared/response"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userUseCase users.UseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase users.UseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetUsers handles GET /api/v1/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUseCase.GetAllUsers(r.Context())
	if err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	userResponses := payload.ToUserResponseList(users)
	response.WriteSuccessResponse(w, "Users retrieved successfully", userResponses)
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req payload.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.userUseCase.CreateUser(r.Context(), user); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Failed to create user", err)
		return
	}

	userResponse := payload.ToUserResponse(user)
	response.WriteCreatedResponse(w, "User created successfully", userResponse)
}

// GetUser handles GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userUseCase.GetUser(r.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, "User not found", err)
		return
	}

	userResponse := payload.ToUserResponse(user)
	response.WriteSuccessResponse(w, "User found", userResponse)
}

// UpdateUser handles PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req payload.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	user := &models.User{
		ID:    uint(id),
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.userUseCase.UpdateUser(r.Context(), user); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Failed to update user", err)
		return
	}

	userResponse := payload.ToUserResponse(user)
	response.WriteSuccessResponse(w, "User updated successfully", userResponse)
}

// DeleteUser handles DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userUseCase.DeleteUser(r.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, "Failed to delete user", err)
		return
	}

	response.WriteSuccessResponse(w, "User deleted successfully", nil)
}
