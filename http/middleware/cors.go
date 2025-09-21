package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
}

// CORS creates a CORS middleware
func CORS(config *CORSConfig) func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins(config.AllowedOrigins),
		handlers.AllowedMethods(config.AllowedMethods),
		handlers.AllowedHeaders(config.AllowedHeaders),
	)
}
