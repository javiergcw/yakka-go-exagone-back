package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/payment_constants/payload"
	"github.com/yakka-backend/internal/features/masters/payment_constants/usecase"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type PaymentConstantHandler struct {
	paymentConstantUsecase usecase.PaymentConstantUsecase
}

func NewPaymentConstantHandler(paymentConstantUsecase usecase.PaymentConstantUsecase) *PaymentConstantHandler {
	return &PaymentConstantHandler{
		paymentConstantUsecase: paymentConstantUsecase,
	}
}

// CreatePaymentConstant creates a new payment constant
func (h *PaymentConstantHandler) CreatePaymentConstant(w http.ResponseWriter, r *http.Request) {
	var req payload.CreatePaymentConstantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create constant
	constant, err := h.paymentConstantUsecase.CreateConstant(r.Context(), req)
	if err != nil {
		if err.Error() == "payment constant with name '"+req.Name+"' already exists" {
			response.WriteError(w, http.StatusConflict, "Payment constant with this name already exists")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to create payment constant")
		return
	}

	// Convert to response
	constantResp := payload.PaymentConstantResponse{
		ID:          constant.ID.String(),
		Name:        constant.Name,
		Value:       constant.Value,
		Description: constant.Description,
		IsActive:    constant.IsActive,
		CreatedAt:   constant.CreatedAt,
		UpdatedAt:   constant.UpdatedAt,
	}

	resp := payload.CreatePaymentConstantResponse{
		Constant: constantResp,
		Message:  "Payment constant created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// GetPaymentConstantByID retrieves a payment constant by ID
func (h *PaymentConstantHandler) GetPaymentConstantByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.WriteError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	constant, err := h.paymentConstantUsecase.GetConstantByID(r.Context(), id)
	if err != nil {
		if err.Error() == "payment constant not found" {
			response.WriteError(w, http.StatusNotFound, "Payment constant not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to get payment constant")
		return
	}

	// Convert to response
	constantResp := payload.PaymentConstantResponse{
		ID:          constant.ID.String(),
		Name:        constant.Name,
		Value:       constant.Value,
		Description: constant.Description,
		IsActive:    constant.IsActive,
		CreatedAt:   constant.CreatedAt,
		UpdatedAt:   constant.UpdatedAt,
	}

	resp := payload.GetPaymentConstantResponse{
		Constant: constantResp,
		Message:  "Payment constant retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetPaymentConstantByName retrieves a payment constant by name
func (h *PaymentConstantHandler) GetPaymentConstantByName(w http.ResponseWriter, r *http.Request) {
	// Extract name from URL path
	name := r.URL.Query().Get("name")
	if name == "" {
		response.WriteError(w, http.StatusBadRequest, "Name parameter is required")
		return
	}

	constant, err := h.paymentConstantUsecase.GetConstantByName(r.Context(), name)
	if err != nil {
		if err.Error() == "payment constant not found" {
			response.WriteError(w, http.StatusNotFound, "Payment constant not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to get payment constant")
		return
	}

	// Convert to response
	constantResp := payload.PaymentConstantResponse{
		ID:          constant.ID.String(),
		Name:        constant.Name,
		Value:       constant.Value,
		Description: constant.Description,
		IsActive:    constant.IsActive,
		CreatedAt:   constant.CreatedAt,
		UpdatedAt:   constant.UpdatedAt,
	}

	resp := payload.GetPaymentConstantResponse{
		Constant: constantResp,
		Message:  "Payment constant retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetAllPaymentConstants retrieves all payment constants
func (h *PaymentConstantHandler) GetAllPaymentConstants(w http.ResponseWriter, r *http.Request) {
	// Simple implementation like licenses
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
		"message": "Payment constants retrieved successfully",
	})
}

// UpdatePaymentConstant updates a payment constant
func (h *PaymentConstantHandler) UpdatePaymentConstant(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.WriteError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req payload.UpdatePaymentConstantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update constant
	constant, err := h.paymentConstantUsecase.UpdateConstant(r.Context(), id, req)
	if err != nil {
		if err.Error() == "payment constant not found" {
			response.WriteError(w, http.StatusNotFound, "Payment constant not found")
			return
		}
		if err.Error() == "payment constant with name '"+*req.Name+"' already exists" {
			response.WriteError(w, http.StatusConflict, "Payment constant with this name already exists")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update payment constant")
		return
	}

	// Convert to response
	constantResp := payload.PaymentConstantResponse{
		ID:          constant.ID.String(),
		Name:        constant.Name,
		Value:       constant.Value,
		Description: constant.Description,
		IsActive:    constant.IsActive,
		CreatedAt:   constant.CreatedAt,
		UpdatedAt:   constant.UpdatedAt,
	}

	resp := payload.UpdatePaymentConstantResponse{
		Constant: constantResp,
		Message:  "Payment constant updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// UpdatePaymentConstantValue updates only the value of a payment constant by name
func (h *PaymentConstantHandler) UpdatePaymentConstantValue(w http.ResponseWriter, r *http.Request) {
	// Extract name from URL path
	name := r.URL.Query().Get("name")
	if name == "" {
		response.WriteError(w, http.StatusBadRequest, "Name parameter is required")
		return
	}

	var req payload.UpdatePaymentConstantValueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update constant value
	constant, err := h.paymentConstantUsecase.UpdateConstantValue(r.Context(), name, req)
	if err != nil {
		if err.Error() == "payment constant not found" {
			response.WriteError(w, http.StatusNotFound, "Payment constant not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update payment constant value")
		return
	}

	// Convert to response
	constantResp := payload.PaymentConstantResponse{
		ID:          constant.ID.String(),
		Name:        constant.Name,
		Value:       constant.Value,
		Description: constant.Description,
		IsActive:    constant.IsActive,
		CreatedAt:   constant.CreatedAt,
		UpdatedAt:   constant.UpdatedAt,
	}

	resp := payload.UpdatePaymentConstantValueResponse{
		Constant: constantResp,
		Message:  "Payment constant value updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// DeletePaymentConstant deletes a payment constant
func (h *PaymentConstantHandler) DeletePaymentConstant(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.WriteError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Delete constant
	err = h.paymentConstantUsecase.DeleteConstant(r.Context(), id)
	if err != nil {
		if err.Error() == "payment constant not found" {
			response.WriteError(w, http.StatusNotFound, "Payment constant not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete payment constant")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Payment constant deleted successfully",
	})
}
