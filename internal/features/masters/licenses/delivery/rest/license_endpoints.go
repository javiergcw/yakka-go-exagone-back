package rest

import (
	"context"
	"net/http"

	licenseRepo "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	"github.com/yakka-backend/internal/infrastructure/database"
	"github.com/yakka-backend/internal/shared/response"
)

type LicenseHandler struct {
	licenseRepository licenseRepo.LicenseRepository
}

func NewLicenseHandler() *LicenseHandler {
	return &LicenseHandler{
		licenseRepository: licenseRepo.NewLicenseRepository(database.DB),
	}
}

// GetLicenses returns all licenses
func (h *LicenseHandler) GetLicenses(w http.ResponseWriter, r *http.Request) {
	licenses, err := h.licenseRepository.GetAll(context.TODO())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to retrieve licenses")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    licenses,
		"message": "Licenses retrieved successfully",
	})
}
