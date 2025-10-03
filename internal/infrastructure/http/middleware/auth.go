package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yakka-backend/internal/infrastructure/database"
	"github.com/yakka-backend/internal/shared/response"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	UserIDKey           ContextKey = "user_id"
	BuilderProfileIDKey ContextKey = "builder_profile_id"
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
			log.Printf("🔐 JWT validation failed: %v", err)
			response.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		log.Printf("🔐 JWT validated successfully for user: %s", claims.UserID)

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		log.Printf("🔍 AuthMiddleware - User ID set in context: %s", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

// BuilderMiddleware validates JWT tokens and ensures user has builder role
func BuilderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First validate JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			response.WriteError(w, http.StatusUnauthorized, "Token required")
			return
		}

		claims, err := validateJWTToken(tokenString)
		if err != nil {
			log.Printf("🔐 JWT validation failed: %v", err)
			response.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Check if user has builder role
		userID := claims.UserID
		var role string
		err = database.DB.Raw("SELECT role FROM users WHERE id = ?", userID).Scan(&role).Error
		if err != nil {
			log.Printf("🔍 Failed to check user role: %v", err)
			response.WriteError(w, http.StatusInternalServerError, "Failed to verify user role")
			return
		}

		if role != "builder" {
			log.Printf("🚫 Access denied: User %s has role %s, required: builder", userID, role)
			response.WriteError(w, http.StatusForbidden, "Access denied: Builder role required")
			return
		}

		log.Printf("🔐 Builder access granted for user: %s", userID)

		// Get builder profile ID for this user
		var builderProfileID string
		err = database.DB.Raw("SELECT id FROM builder_profiles WHERE user_id = ?", userID).Scan(&builderProfileID).Error
		if err != nil {
			log.Printf("🔍 Failed to get builder profile ID: %v", err)
			response.WriteError(w, http.StatusInternalServerError, "Failed to get builder profile")
			return
		}

		if builderProfileID == "" {
			log.Printf("🚫 No builder profile found for user: %s", userID)
			response.WriteError(w, http.StatusForbidden, "Builder profile not found")
			return
		}

		// Add user ID and builder profile ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, BuilderProfileIDKey, builderProfileID)

		log.Printf("🔍 Context values set - UserID: %s, BuilderProfileID: %s", userID, builderProfileID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LabourMiddleware validates JWT tokens and ensures user has labour role
func LabourMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First validate JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			response.WriteError(w, http.StatusUnauthorized, "Token required")
			return
		}

		claims, err := validateJWTToken(tokenString)
		if err != nil {
			log.Printf("🔐 JWT validation failed: %v", err)
			response.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Check if user has labour role
		userID := claims.UserID
		var role string
		err = database.DB.Raw("SELECT role FROM users WHERE id = ?", userID).Scan(&role).Error
		if err != nil {
			log.Printf("🔍 Failed to check user role: %v", err)
			response.WriteError(w, http.StatusInternalServerError, "Failed to verify user role")
			return
		}

		if role != "labour" {
			log.Printf("🚫 Access denied: User %s has role %s, required: labour", userID, role)
			response.WriteError(w, http.StatusForbidden, "Access denied: Labour role required")
			return
		}

		log.Printf("🔐 Labour access granted for user: %s", userID)

		// Set context values for downstream handlers
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		log.Printf("🔍 Context values set - UserID: %s", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
