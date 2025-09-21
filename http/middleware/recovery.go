package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				// Set content type
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				// Return error response
				errorResponse := map[string]interface{}{
					"error":   "Internal Server Error",
					"message": "Something went wrong",
				}

				json.NewEncoder(w).Encode(errorResponse)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
