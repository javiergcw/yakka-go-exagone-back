package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/job_assignments/models"
	"github.com/yakka-backend/internal/features/job_assignments/payload"
	"github.com/yakka-backend/internal/features/job_assignments/usecase"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type JobAssignmentHandler struct {
	assignmentUsecase usecase.JobAssignmentUsecase
}

func NewJobAssignmentHandler(assignmentUsecase usecase.JobAssignmentUsecase) *JobAssignmentHandler {
	return &JobAssignmentHandler{
		assignmentUsecase: assignmentUsecase,
	}
}

// CreateAssignment creates a new job assignment
func (h *JobAssignmentHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var req payload.CreateJobAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create assignment
	assignment, err := h.assignmentUsecase.CreateAssignment(r.Context(), req)
	if err != nil {
		if err.Error() == "assignment already exists for this application" {
			response.WriteError(w, http.StatusConflict, "Assignment already exists for this application")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to create assignment")
		return
	}

	// Convert to response
	assignmentResp := payload.JobAssignmentResponse{
		ID:            assignment.ID.String(),
		JobID:         assignment.JobID.String(),
		LabourUserID:  assignment.LabourUserID.String(),
		ApplicationID: assignment.ApplicationID.String(),
		StartDate:     assignment.StartDate,
		EndDate:       assignment.EndDate,
		Status:        assignment.Status,
		CreatedAt:     assignment.CreatedAt,
		UpdatedAt:     assignment.UpdatedAt,
	}

	resp := payload.CreateJobAssignmentResponse{
		Assignment: assignmentResp,
		Message:    "Assignment created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// GetAssignmentByID retrieves a job assignment by ID
func (h *JobAssignmentHandler) GetAssignmentByID(w http.ResponseWriter, r *http.Request) {
	// Get assignment ID from URL
	vars := r.URL.Query()
	assignmentIDStr := vars.Get("id")
	if assignmentIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	// Get assignment
	assignment, err := h.assignmentUsecase.GetAssignmentByID(r.Context(), assignmentID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Assignment not found")
		return
	}

	// Convert to response
	assignmentResp := payload.JobAssignmentResponse{
		ID:            assignment.ID.String(),
		JobID:         assignment.JobID.String(),
		LabourUserID:  assignment.LabourUserID.String(),
		ApplicationID: assignment.ApplicationID.String(),
		StartDate:     assignment.StartDate,
		EndDate:       assignment.EndDate,
		Status:        assignment.Status,
		CreatedAt:     assignment.CreatedAt,
		UpdatedAt:     assignment.UpdatedAt,
	}

	response.WriteJSON(w, http.StatusOK, assignmentResp)
}

// UpdateAssignment updates an existing job assignment
func (h *JobAssignmentHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	// Get assignment ID from URL
	vars := r.URL.Query()
	assignmentIDStr := vars.Get("id")
	if assignmentIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var req payload.UpdateJobAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update assignment
	assignment, err := h.assignmentUsecase.UpdateAssignment(r.Context(), assignmentID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update assignment")
		return
	}

	// Convert to response
	assignmentResp := payload.JobAssignmentResponse{
		ID:            assignment.ID.String(),
		JobID:         assignment.JobID.String(),
		LabourUserID:  assignment.LabourUserID.String(),
		ApplicationID: assignment.ApplicationID.String(),
		StartDate:     assignment.StartDate,
		EndDate:       assignment.EndDate,
		Status:        assignment.Status,
		CreatedAt:     assignment.CreatedAt,
		UpdatedAt:     assignment.UpdatedAt,
	}

	resp := payload.UpdateJobAssignmentResponse{
		Assignment: assignmentResp,
		Message:    "Assignment updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// DeleteAssignment deletes a job assignment
func (h *JobAssignmentHandler) DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	// Get assignment ID from URL
	vars := r.URL.Query()
	assignmentIDStr := vars.Get("id")
	if assignmentIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	// Delete assignment
	if err := h.assignmentUsecase.DeleteAssignment(r.Context(), assignmentID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete assignment")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Assignment deleted successfully"})
}

// GetAssignments retrieves job assignments with filters
func (h *JobAssignmentHandler) GetAssignments(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := payload.GetJobAssignmentsRequest{
		JobID:         getStringParam(r, "job_id"),
		LabourUserID:  getStringParam(r, "labour_user_id"),
		ApplicationID: getStringParam(r, "application_id"),
		Status:        getStringParam(r, "status"),
		Page:          getIntParam(r, "page", 1),
		Limit:         getIntParam(r, "limit", 20),
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get assignments
	assignments, total, err := h.assignmentUsecase.GetAssignmentsWithFilters(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get assignments")
		return
	}

	// Convert to response
	assignmentResponses := make([]payload.JobAssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		assignmentResponses[i] = payload.JobAssignmentResponse{
			ID:            assignment.ID.String(),
			JobID:         assignment.JobID.String(),
			LabourUserID:  assignment.LabourUserID.String(),
			ApplicationID: assignment.ApplicationID.String(),
			StartDate:     assignment.StartDate,
			EndDate:       assignment.EndDate,
			Status:        assignment.Status,
			CreatedAt:     assignment.CreatedAt,
			UpdatedAt:     assignment.UpdatedAt,
		}
	}

	totalPages := calculateTotalPages(total, req.Limit)

	resp := payload.GetJobAssignmentsResponse{
		Assignments: assignmentResponses,
		Total:       total,
		Page:        req.Page,
		Limit:       req.Limit,
		TotalPages:  totalPages,
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// UpdateAssignmentStatus updates the status of a job assignment
func (h *JobAssignmentHandler) UpdateAssignmentStatus(w http.ResponseWriter, r *http.Request) {
	// Get assignment ID from URL
	vars := r.URL.Query()
	assignmentIDStr := vars.Get("id")
	if assignmentIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=ACTIVE COMPLETED CANCELLED"`
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

	status := models.AssignmentStatus(req.Status)

	// Update status
	assignment, err := h.assignmentUsecase.UpdateAssignmentStatus(r.Context(), assignmentID, status)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update assignment status")
		return
	}

	// Convert to response
	assignmentResp := payload.JobAssignmentResponse{
		ID:            assignment.ID.String(),
		JobID:         assignment.JobID.String(),
		LabourUserID:  assignment.LabourUserID.String(),
		ApplicationID: assignment.ApplicationID.String(),
		StartDate:     assignment.StartDate,
		EndDate:       assignment.EndDate,
		Status:        assignment.Status,
		CreatedAt:     assignment.CreatedAt,
		UpdatedAt:     assignment.UpdatedAt,
	}

	response.WriteJSON(w, http.StatusOK, assignmentResp)
}

// CompleteAssignment completes a job assignment
func (h *JobAssignmentHandler) CompleteAssignment(w http.ResponseWriter, r *http.Request) {
	// Get assignment ID from URL
	vars := r.URL.Query()
	assignmentIDStr := vars.Get("id")
	if assignmentIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var req payload.CompleteAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Complete assignment
	assignment, err := h.assignmentUsecase.CompleteAssignment(r.Context(), assignmentID, req)
	if err != nil {
		if err.Error() == "assignment already completed" {
			response.WriteError(w, http.StatusConflict, "Assignment already completed")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to complete assignment")
		return
	}

	// Convert to response
	assignmentResp := payload.JobAssignmentResponse{
		ID:            assignment.ID.String(),
		JobID:         assignment.JobID.String(),
		LabourUserID:  assignment.LabourUserID.String(),
		ApplicationID: assignment.ApplicationID.String(),
		StartDate:     assignment.StartDate,
		EndDate:       assignment.EndDate,
		Status:        assignment.Status,
		CreatedAt:     assignment.CreatedAt,
		UpdatedAt:     assignment.UpdatedAt,
	}

	resp := payload.CompleteAssignmentResponse{
		Assignment: assignmentResp,
		Message:    "Assignment completed successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// CancelAssignment cancels a job assignment
func (h *JobAssignmentHandler) CancelAssignment(w http.ResponseWriter, r *http.Request) {
	// Get assignment ID from URL
	vars := r.URL.Query()
	assignmentIDStr := vars.Get("id")
	if assignmentIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var req payload.CancelAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Cancel assignment
	assignment, err := h.assignmentUsecase.CancelAssignment(r.Context(), assignmentID, req)
	if err != nil {
		if err.Error() == "assignment already cancelled" {
			response.WriteError(w, http.StatusConflict, "Assignment already cancelled")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to cancel assignment")
		return
	}

	// Convert to response
	assignmentResp := payload.JobAssignmentResponse{
		ID:            assignment.ID.String(),
		JobID:         assignment.JobID.String(),
		LabourUserID:  assignment.LabourUserID.String(),
		ApplicationID: assignment.ApplicationID.String(),
		StartDate:     assignment.StartDate,
		EndDate:       assignment.EndDate,
		Status:        assignment.Status,
		CreatedAt:     assignment.CreatedAt,
		UpdatedAt:     assignment.UpdatedAt,
	}

	resp := payload.CancelAssignmentResponse{
		Assignment: assignmentResp,
		Message:    "Assignment cancelled successfully",
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
