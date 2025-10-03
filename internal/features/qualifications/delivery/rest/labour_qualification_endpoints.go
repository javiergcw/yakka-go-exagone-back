package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/qualifications/entity/database"
	"github.com/yakka-backend/internal/features/qualifications/models"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
)

// LabourQualificationHandler handles labour qualification-related HTTP requests
type LabourQualificationHandler struct {
	labourQualificationRepo database.LabourProfileQualificationRepository
	qualificationRepo       database.QualificationRepository
}

// NewLabourQualificationHandler creates a new labour qualification handler
func NewLabourQualificationHandler(
	labourQualificationRepo database.LabourProfileQualificationRepository,
	qualificationRepo database.QualificationRepository,
) *LabourQualificationHandler {
	return &LabourQualificationHandler{
		labourQualificationRepo: labourQualificationRepo,
		qualificationRepo:       qualificationRepo,
	}
}

// LabourQualificationItem represents a qualification assigned to a labour profile
type LabourQualificationItem struct {
	ID              string     `json:"id"`
	QualificationID string     `json:"qualification_id"`
	Title           string     `json:"title"`
	Organization    string     `json:"organization,omitempty"`
	Country         string     `json:"country,omitempty"`
	Sport           string     `json:"sport"`
	DateObtained    *time.Time `json:"date_obtained,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	Status          string     `json:"status"`
}

// GetLabourQualificationsResponse represents the response for getting labour qualifications
type GetLabourQualificationsResponse struct {
	Qualifications []LabourQualificationItem `json:"qualifications"`
	Total          int                       `json:"total"`
	Message        string                    `json:"message"`
}

// CreateLabourQualificationsRequest represents the request to create labour qualifications
type CreateLabourQualificationsRequest struct {
	Qualifications []struct {
		QualificationID string     `json:"qualification_id" validate:"required,uuid"`
		DateObtained    *time.Time `json:"date_obtained,omitempty"`
		ExpiresAt       *time.Time `json:"expires_at,omitempty"`
		Status          string     `json:"status,omitempty"`
	} `json:"qualifications" validate:"required"`
}

// UpdateLabourQualificationsRequest represents the request to update labour qualifications
type UpdateLabourQualificationsRequest struct {
	Qualifications []struct {
		QualificationID string     `json:"qualification_id" validate:"required,uuid"`
		DateObtained    *time.Time `json:"date_obtained,omitempty"`
		ExpiresAt       *time.Time `json:"expires_at,omitempty"`
		Status          string     `json:"status,omitempty"`
	} `json:"qualifications" validate:"required"`
}

// CreateLabourQualificationsResponse represents the response for creating labour qualifications
type CreateLabourQualificationsResponse struct {
	Qualifications []LabourQualificationItem `json:"qualifications"`
	Message        string                    `json:"message"`
}

// UpdateLabourQualificationsResponse represents the response for updating labour qualifications
type UpdateLabourQualificationsResponse struct {
	Qualifications []LabourQualificationItem `json:"qualifications"`
	Message        string                    `json:"message"`
}

// GetLabourQualifications handles GET /api/v1/labour/qualifications
func (h *LabourQualificationHandler) GetLabourQualifications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get user ID from context (set by LabourMiddleware)
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get labour profile qualifications
	labourQualifications, err := h.labourQualificationRepo.GetByLabourProfileID(ctx, userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get labour qualifications")
		return
	}

	// Convert to response format
	qualificationItems := make([]LabourQualificationItem, len(labourQualifications))
	for i, lq := range labourQualifications {
		qualificationItems[i] = LabourQualificationItem{
			ID:              lq.ID.String(),
			QualificationID: lq.QualificationID.String(),
			Title:           getQualificationTitle(lq.Qualification),
			Organization:    getQualificationOrganization(lq.Qualification),
			Country:         getQualificationCountry(lq.Qualification),
			Sport:           getQualificationSport(lq.Qualification),
			DateObtained:    lq.DateObtained,
			ExpiresAt:       lq.ExpiresAt,
			Status:          lq.Status,
		}
	}

	resp := GetLabourQualificationsResponse{
		Qualifications: qualificationItems,
		Total:          len(qualificationItems),
		Message:        "Labour qualifications retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// CreateLabourQualifications handles POST /api/v1/labour/qualifications
func (h *LabourQualificationHandler) CreateLabourQualifications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get user ID from context (set by LabourMiddleware)
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Parse request body
	var req CreateLabourQualificationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if len(req.Qualifications) == 0 {
		response.WriteError(w, http.StatusBadRequest, "At least one qualification is required")
		return
	}

	// Validate qualification IDs exist
	for _, q := range req.Qualifications {
		qualificationID, err := uuid.Parse(q.QualificationID)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "Invalid qualification ID: "+q.QualificationID)
			return
		}

		// Check if qualification exists
		qualification, err := h.qualificationRepo.GetByID(ctx, qualificationID)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "Qualification not found: "+q.QualificationID)
			return
		}

		if qualification == nil {
			response.WriteError(w, http.StatusBadRequest, "Qualification not found: "+q.QualificationID)
			return
		}
	}

	// Create new qualifications
	var newQualifications []LabourQualificationItem
	for _, q := range req.Qualifications {
		qualificationID, _ := uuid.Parse(q.QualificationID)

		labourQualification := &models.LabourProfileQualification{
			LabourProfileID: userID,
			QualificationID: qualificationID,
			DateObtained:    q.DateObtained,
			ExpiresAt:       q.ExpiresAt,
			Status:          q.Status,
		}

		if labourQualification.Status == "" {
			labourQualification.Status = "valid"
		}

		if err := h.labourQualificationRepo.Create(ctx, labourQualification); err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Failed to create qualification")
			return
		}

		// Get the created qualification with details
		createdQualification, err := h.labourQualificationRepo.GetByID(ctx, labourQualification.ID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve created qualification")
			return
		}

		newQualifications = append(newQualifications, LabourQualificationItem{
			ID:              createdQualification.ID.String(),
			QualificationID: createdQualification.QualificationID.String(),
			Title:           getQualificationTitle(createdQualification.Qualification),
			Organization:    getQualificationOrganization(createdQualification.Qualification),
			Country:         getQualificationCountry(createdQualification.Qualification),
			Sport:           getQualificationSport(createdQualification.Qualification),
			DateObtained:    createdQualification.DateObtained,
			ExpiresAt:       createdQualification.ExpiresAt,
			Status:          createdQualification.Status,
		})
	}

	resp := CreateLabourQualificationsResponse{
		Qualifications: newQualifications,
		Message:        "Labour qualifications created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// UpdateLabourQualifications handles POST /api/v1/labour/qualifications
func (h *LabourQualificationHandler) UpdateLabourQualifications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get user ID from context (set by LabourMiddleware)
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Parse request body
	var req UpdateLabourQualificationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if len(req.Qualifications) == 0 {
		response.WriteError(w, http.StatusBadRequest, "At least one qualification is required")
		return
	}

	// Validate qualification IDs exist
	for _, q := range req.Qualifications {
		qualificationID, err := uuid.Parse(q.QualificationID)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "Invalid qualification ID: "+q.QualificationID)
			return
		}

		// Check if qualification exists
		qualification, err := h.qualificationRepo.GetByID(ctx, qualificationID)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "Qualification not found: "+q.QualificationID)
			return
		}

		if qualification == nil {
			response.WriteError(w, http.StatusBadRequest, "Qualification not found: "+q.QualificationID)
			return
		}
	}

	// Delete existing qualifications for this labour profile
	if err := h.labourQualificationRepo.DeleteByLabourProfileID(ctx, userID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete existing qualifications")
		return
	}

	// Create new qualifications
	var newQualifications []LabourQualificationItem
	for _, q := range req.Qualifications {
		qualificationID, _ := uuid.Parse(q.QualificationID)

		labourQualification := &models.LabourProfileQualification{
			LabourProfileID: userID,
			QualificationID: qualificationID,
			DateObtained:    q.DateObtained,
			ExpiresAt:       q.ExpiresAt,
			Status:          q.Status,
		}

		if labourQualification.Status == "" {
			labourQualification.Status = "valid"
		}

		if err := h.labourQualificationRepo.Create(ctx, labourQualification); err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Failed to create qualification")
			return
		}

		// Get the created qualification with details
		createdQualification, err := h.labourQualificationRepo.GetByID(ctx, labourQualification.ID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve created qualification")
			return
		}

		newQualifications = append(newQualifications, LabourQualificationItem{
			ID:              createdQualification.ID.String(),
			QualificationID: createdQualification.QualificationID.String(),
			Title:           getQualificationTitle(createdQualification.Qualification),
			Organization:    getQualificationOrganization(createdQualification.Qualification),
			Country:         getQualificationCountry(createdQualification.Qualification),
			Sport:           getQualificationSport(createdQualification.Qualification),
			DateObtained:    createdQualification.DateObtained,
			ExpiresAt:       createdQualification.ExpiresAt,
			Status:          createdQualification.Status,
		})
	}

	resp := UpdateLabourQualificationsResponse{
		Qualifications: newQualifications,
		Message:        "Labour qualifications updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// Helper functions
func getQualificationTitle(qualification *models.Qualification) string {
	if qualification == nil {
		return ""
	}
	return qualification.Title
}

func getQualificationOrganization(qualification *models.Qualification) string {
	if qualification == nil || qualification.Organization == nil {
		return ""
	}
	return *qualification.Organization
}

func getQualificationCountry(qualification *models.Qualification) string {
	if qualification == nil || qualification.Country == nil {
		return ""
	}
	return *qualification.Country
}

func getQualificationSport(qualification *models.Qualification) string {
	if qualification == nil || qualification.Sport == nil {
		return ""
	}
	return qualification.Sport.Name
}
