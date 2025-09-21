package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Version   string                 `json:"version"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccessResponse writes a success response
func WriteSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	WriteJSON(w, http.StatusOK, response)
}

// WriteCreatedResponse writes a created response
func WriteCreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	WriteJSON(w, http.StatusCreated, response)
}

// WriteErrorResponse writes an error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	WriteJSON(w, statusCode, response)
}

// WriteBadRequestResponse writes a bad request response
func WriteBadRequestResponse(w http.ResponseWriter, message string, err error) {
	WriteErrorResponse(w, http.StatusBadRequest, message, err)
}

// WriteNotFoundResponse writes a not found response
func WriteNotFoundResponse(w http.ResponseWriter, message string, err error) {
	WriteErrorResponse(w, http.StatusNotFound, message, err)
}

// WriteInternalServerErrorResponse writes an internal server error response
func WriteInternalServerErrorResponse(w http.ResponseWriter, message string, err error) {
	WriteErrorResponse(w, http.StatusInternalServerError, message, err)
}

// WriteError writes a simple error response
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	response := Response{
		Success: false,
		Message: message,
		Error:   message,
	}
	WriteJSON(w, statusCode, response)
}
