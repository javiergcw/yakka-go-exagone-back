package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yakka-backend/internal/shared/response"
)

// JWTSecret should be loaded from environment variables in production
var JWTSecret = []byte("your-secret-key-change-in-production")

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check and public auth endpoints
		if isPublicEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			response.WriteError(w, http.StatusUnauthorized, "Token required")
			return
		}

		// Parse and validate JWT token
		claims, err := validateJWTToken(tokenString)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isPublicEndpoint checks if the endpoint is public (doesn't require authentication)
func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/auth/email/request-verification",
		"/api/v1/auth/email/verify",
		"/api/v1/auth/password/request-reset",
		"/api/v1/auth/password/reset",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath {
			return true
		}
	}
	return false
}

// validateJWTToken validates and parses a JWT token
func validateJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

// GenerateJWTToken generates a new JWT token for a user
func GenerateJWTToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
