package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/yakka-backend/internal/features/jobsites/payload"
	"github.com/yakka-backend/internal/features/jobsites/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

// JobsiteHandler handles jobsite HTTP requests
type JobsiteHandler struct {
	jobsiteUsecase usecase.JobsiteUsecase
}

// NewJobsiteHandler creates a new instance of JobsiteHandler
func NewJobsiteHandler(jobsiteUsecase usecase.JobsiteUsecase) *JobsiteHandler {
	return &JobsiteHandler{
		jobsiteUsecase: jobsiteUsecase,
	}
}

// CreateJobsite creates a new jobsite for the authenticated builder
func (h *JobsiteHandler) CreateJobsite(w http.ResponseWriter, r *http.Request) {
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

	var req payload.CreateJobsiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set the builder_id to the authenticated user's ID
	req.BuilderID = userID.String()

	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.jobsiteUsecase.CreateJobsite(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, result)
}

// GetJobsitesByBuilder retrieves all jobsites for the authenticated builder
func (h *JobsiteHandler) GetJobsitesByBuilder(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.jobsiteUsecase.GetJobsitesByBuilderID(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

// GetJobsiteByID retrieves a specific jobsite by ID (only if it belongs to the authenticated builder)
func (h *JobsiteHandler) GetJobsiteByID(w http.ResponseWriter, r *http.Request) {
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

	// Get jobsite ID from URL path
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		response.WriteError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid jobsite ID")
		return
	}

	// Get the jobsite
	jobsite, err := h.jobsiteUsecase.GetJobsiteByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	// Verify that the jobsite belongs to the authenticated builder
	if jobsite.BuilderID != userID {
		response.WriteError(w, http.StatusForbidden, "Access denied: This jobsite does not belong to you")
		return
	}

	response.WriteJSON(w, http.StatusOK, jobsite)
}

// UpdateJobsite updates a specific jobsite (only if it belongs to the authenticated builder)
func (h *JobsiteHandler) UpdateJobsite(w http.ResponseWriter, r *http.Request) {
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

	// Get jobsite ID from URL path
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		response.WriteError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid jobsite ID")
		return
	}

	// First, verify that the jobsite belongs to the authenticated builder
	jobsite, err := h.jobsiteUsecase.GetJobsiteByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if jobsite.BuilderID != userID {
		response.WriteError(w, http.StatusForbidden, "Access denied: This jobsite does not belong to you")
		return
	}

	var req payload.UpdateJobsiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.jobsiteUsecase.UpdateJobsite(r.Context(), id, &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

// DeleteJobsite deletes a specific jobsite (only if it belongs to the authenticated builder)
func (h *JobsiteHandler) DeleteJobsite(w http.ResponseWriter, r *http.Request) {
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

	// Get jobsite ID from URL path
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		response.WriteError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid jobsite ID")
		return
	}

	// First, verify that the jobsite belongs to the authenticated builder
	jobsite, err := h.jobsiteUsecase.GetJobsiteByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if jobsite.BuilderID != userID {
		response.WriteError(w, http.StatusForbidden, "Access denied: This jobsite does not belong to you")
		return
	}

	result, err := h.jobsiteUsecase.DeleteJobsite(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}
