package middleware

import (
	"net/http"
	"strings"

	"github.com/yakka-backend/internal/shared/response"
)

// LicenseMiddleware validates license for endpoints
func LicenseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip license validation for health check and public endpoints
		if isPublicEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get License header
		licenseHeader := r.Header.Get("X-License")
		if licenseHeader == "" {
			response.WriteError(w, http.StatusUnauthorized, "License header required")
			return
		}

		// Validate license
		if !isValidLicense(licenseHeader) {
			response.WriteError(w, http.StatusUnauthorized, "Invalid license")
			return
		}

		next.ServeHTTP(w, r)
	})
}


// isValidLicense validates the license key
func isValidLicense(license string) bool {
	license = strings.TrimSpace(license)
	
	// Check for the specific production license
	validLicense := "YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"
	
	return license == validLicense
}
