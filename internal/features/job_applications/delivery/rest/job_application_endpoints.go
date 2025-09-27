package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_applications/models"
	"github.com/yakka-backend/internal/features/job_applications/payload"
	"github.com/yakka-backend/internal/features/job_applications/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type JobApplicationHandler struct {
	applicationUsecase usecase.JobApplicationUsecase
}

func NewJobApplicationHandler(applicationUsecase usecase.JobApplicationUsecase) *JobApplicationHandler {
	return &JobApplicationHandler{
		applicationUsecase: applicationUsecase,
	}
}

// CreateApplication creates a new job application
func (h *JobApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
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

	var req payload.CreateJobApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create application
	application, err := h.applicationUsecase.CreateApplication(r.Context(), userID, req)
	if err != nil {
		if err.Error() == "application already exists for this job" {
			response.WriteError(w, http.StatusConflict, "Application already exists for this job")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to create application")
		return
	}

	// Convert to response
	applicationResp := payload.JobApplicationResponse{
		ID:           application.ID.String(),
		JobID:        application.JobID.String(),
		LabourUserID: application.LabourUserID.String(),
		Status:       application.Status,
		CoverLetter:  application.CoverLetter,
		ExpectedRate: application.ExpectedRate,
		ResumeURL:    application.ResumeURL,
		CreatedAt:    application.CreatedAt,
		UpdatedAt:    application.UpdatedAt,
		WithdrawnAt:  application.WithdrawnAt,
	}

	resp := payload.CreateJobApplicationResponse{
		Application: applicationResp,
		Message:     "Application created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// GetApplicationByID retrieves a job application by ID
func (h *JobApplicationHandler) GetApplicationByID(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	vars := r.URL.Query()
	applicationIDStr := vars.Get("id")
	if applicationIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}

	// Get application
	application, err := h.applicationUsecase.GetApplicationByID(r.Context(), applicationID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Application not found")
		return
	}

	// Convert to response
	applicationResp := payload.JobApplicationResponse{
		ID:           application.ID.String(),
		JobID:        application.JobID.String(),
		LabourUserID: application.LabourUserID.String(),
		Status:       application.Status,
		CoverLetter:  application.CoverLetter,
		ExpectedRate: application.ExpectedRate,
		ResumeURL:    application.ResumeURL,
		CreatedAt:    application.CreatedAt,
		UpdatedAt:    application.UpdatedAt,
		WithdrawnAt:  application.WithdrawnAt,
	}

	response.WriteJSON(w, http.StatusOK, applicationResp)
}

// UpdateApplication updates an existing job application
func (h *JobApplicationHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	vars := r.URL.Query()
	applicationIDStr := vars.Get("id")
	if applicationIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}

	var req payload.UpdateJobApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update application
	application, err := h.applicationUsecase.UpdateApplication(r.Context(), applicationID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update application")
		return
	}

	// Convert to response
	applicationResp := payload.JobApplicationResponse{
		ID:           application.ID.String(),
		JobID:        application.JobID.String(),
		LabourUserID: application.LabourUserID.String(),
		Status:       application.Status,
		CoverLetter:  application.CoverLetter,
		ExpectedRate: application.ExpectedRate,
		ResumeURL:    application.ResumeURL,
		CreatedAt:    application.CreatedAt,
		UpdatedAt:    application.UpdatedAt,
		WithdrawnAt:  application.WithdrawnAt,
	}

	resp := payload.UpdateJobApplicationResponse{
		Application: applicationResp,
		Message:     "Application updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// DeleteApplication deletes a job application
func (h *JobApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	vars := r.URL.Query()
	applicationIDStr := vars.Get("id")
	if applicationIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}

	// Delete application
	if err := h.applicationUsecase.DeleteApplication(r.Context(), applicationID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete application")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Application deleted successfully"})
}

// GetApplications retrieves job applications with filters
func (h *JobApplicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := payload.GetJobApplicationsRequest{
		JobID:        getStringParam(r, "job_id"),
		LabourUserID: getStringParam(r, "labour_user_id"),
		Status:       getStringParam(r, "status"),
		Page:         getIntParam(r, "page", 1),
		Limit:        getIntParam(r, "limit", 20),
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get applications
	applications, total, err := h.applicationUsecase.GetApplicationsWithFilters(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get applications")
		return
	}

	// Convert to response
	applicationResponses := make([]payload.JobApplicationResponse, len(applications))
	for i, app := range applications {
		applicationResponses[i] = payload.JobApplicationResponse{
			ID:           app.ID.String(),
			JobID:        app.JobID.String(),
			LabourUserID: app.LabourUserID.String(),
			Status:       app.Status,
			CoverLetter:  app.CoverLetter,
			ExpectedRate: app.ExpectedRate,
			ResumeURL:    app.ResumeURL,
			CreatedAt:    app.CreatedAt,
			UpdatedAt:    app.UpdatedAt,
			WithdrawnAt:  app.WithdrawnAt,
		}
	}

	totalPages := calculateTotalPages(total, req.Limit)

	resp := payload.GetJobApplicationsResponse{
		Applications: applicationResponses,
		Total:        total,
		Page:         req.Page,
		Limit:        req.Limit,
		TotalPages:   totalPages,
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// UpdateApplicationStatus updates the status of a job application
func (h *JobApplicationHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	vars := r.URL.Query()
	applicationIDStr := vars.Get("id")
	if applicationIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=APPLIED REVIEWED ACCEPTED REJECTED WITHDRAWN"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	status := models.ApplicationStatus(req.Status)

	// Update status
	application, err := h.applicationUsecase.UpdateApplicationStatus(r.Context(), applicationID, status)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update application status")
		return
	}

	// Convert to response
	applicationResp := payload.JobApplicationResponse{
		ID:           application.ID.String(),
		JobID:        application.JobID.String(),
		LabourUserID: application.LabourUserID.String(),
		Status:       application.Status,
		CoverLetter:  application.CoverLetter,
		ExpectedRate: application.ExpectedRate,
		ResumeURL:    application.ResumeURL,
		CreatedAt:    application.CreatedAt,
		UpdatedAt:    application.UpdatedAt,
		WithdrawnAt:  application.WithdrawnAt,
	}

	response.WriteJSON(w, http.StatusOK, applicationResp)
}

// WithdrawApplication withdraws a job application
func (h *JobApplicationHandler) WithdrawApplication(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	vars := r.URL.Query()
	applicationIDStr := vars.Get("id")
	if applicationIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}

	// Withdraw application
	if err := h.applicationUsecase.WithdrawApplication(r.Context(), applicationID); err != nil {
		if err.Error() == "application already withdrawn" {
			response.WriteError(w, http.StatusConflict, "Application already withdrawn")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to withdraw application")
		return
	}

	resp := payload.WithdrawApplicationResponse{
		Message: "Application withdrawn successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// Helper functions
func getStringParam(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func getIntParam(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func calculateTotalPages(total int64, limit int) int {
	if limit == 0 {
		return 0
	}
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}
	return pages
}
