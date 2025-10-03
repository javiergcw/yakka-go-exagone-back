package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/builder_profiles/payload"
	"github.com/yakka-backend/internal/features/builder_profiles/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

// CompanyHandler handles company endpoints
type CompanyHandler struct {
	companyUsecase usecase.CompanyUsecase
}

func NewCompanyHandler(companyUsecase usecase.CompanyUsecase) *CompanyHandler {
	return &CompanyHandler{
		companyUsecase: companyUsecase,
	}
}

// CreateCompany creates a new company
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req payload.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create company
	company, err := h.companyUsecase.CreateCompany(r.Context(), req)
	if err != nil {
		// Check for specific error types
		if err.Error() == "company with this name already exists" {
			response.WriteError(w, http.StatusConflict, "Company with this name already exists")
			return
		}

		// Generic error for other cases
		response.WriteError(w, http.StatusInternalServerError, "Failed to create company")
		return
	}

	// Convert to response
	companyResp := payload.CompanyResponse{
		ID:          company.ID.String(),
		Name:        company.Name,
		Description: company.Description,
		Website:     company.Website,
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
	}

	resp := payload.CreateCompanyResponse{
		Company: companyResp,
		Message: "Company created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// GetCompanies retrieves all companies
func (h *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	// Get all companies
	companies, err := h.companyUsecase.GetAllCompanies(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get companies")
		return
	}

	// Convert to response
	var companiesResp []payload.CompanyResponse
	for _, company := range companies {
		companyResp := payload.CompanyResponse{
			ID:          company.ID.String(),
			Name:        company.Name,
			Description: company.Description,
			Website:     company.Website,
			CreatedAt:   company.CreatedAt,
			UpdatedAt:   company.UpdatedAt,
		}
		companiesResp = append(companiesResp, companyResp)
	}

	resp := payload.GetCompaniesResponse{
		Companies: companiesResp,
		Total:     len(companiesResp),
		Message:   "Companies retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// AssignCompany assigns a company to a builder
func (h *CompanyHandler) AssignCompany(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by BuilderMiddleware)
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

	var req payload.AssignCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Assign company to builder
	builderProfile, err := h.companyUsecase.AssignCompanyToBuilder(r.Context(), userID, req)
	if err != nil {
		// Check for specific error types
		if err.Error() == "company not found" {
			response.WriteError(w, http.StatusBadRequest, "Company not found")
			return
		}
		if err.Error() == "invalid company_id" {
			response.WriteError(w, http.StatusBadRequest, "Invalid company ID")
			return
		}

		// Generic error for other cases
		response.WriteError(w, http.StatusInternalServerError, "Failed to assign company")
		return
	}

	// Convert to response
	profileResp := payload.BuilderProfileResponse{
		ID:          builderProfile.ID.String(),
		UserID:      builderProfile.UserID.String(),
		DisplayName: getStringValue(builderProfile.DisplayName),
		Location:    getStringValue(builderProfile.Location),
		Bio:         builderProfile.Bio,
		CreatedAt:   builderProfile.CreatedAt,
		UpdatedAt:   builderProfile.UpdatedAt,
	}

	// Add company information if available
	if builderProfile.CompanyID != nil {
		companyIDStr := builderProfile.CompanyID.String()
		profileResp.CompanyID = &companyIDStr
	}

	resp := payload.AssignCompanyResponse{
		BuilderProfile: profileResp,
		Message:        "Company assigned successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
