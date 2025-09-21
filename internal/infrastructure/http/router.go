package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	auth_rest "github.com/yakka-backend/internal/features/auth/delivery/rest"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
)

// Router sets up the HTTP routes
type Router struct {
	authHandler     *auth_rest.AuthHandler
	sessionHandler  *auth_rest.SessionHandler
	passwordHandler *auth_rest.PasswordHandler
	emailHandler    *auth_rest.EmailHandler
}

// NewRouter creates a new router
func NewRouter(
	authHandler *auth_rest.AuthHandler,
	sessionHandler *auth_rest.SessionHandler,
	passwordHandler *auth_rest.PasswordHandler,
	emailHandler *auth_rest.EmailHandler,
) *Router {
	return &Router{
		authHandler:     authHandler,
		sessionHandler:  sessionHandler,
		passwordHandler: passwordHandler,
		emailHandler:    emailHandler,
	}
}

// SetupRoutes configures all the routes
func (r *Router) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", r.healthCheck).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Auth endpoints (public)
	api.HandleFunc("/auth/register", r.authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", r.authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", r.sessionHandler.RefreshToken).Methods("POST")
	api.HandleFunc("/auth/password/reset", r.passwordHandler.RequestPasswordReset).Methods("POST")
	api.HandleFunc("/auth/password/reset/confirm", r.passwordHandler.ResetPassword).Methods("POST")
	api.HandleFunc("/auth/email/verify", r.emailHandler.VerifyEmail).Methods("POST")

	// Protected auth endpoints
	api.HandleFunc("/auth/profile", r.authHandler.GetProfile).Methods("GET")
	api.HandleFunc("/auth/profile", r.authHandler.UpdateProfile).Methods("PUT")
	api.HandleFunc("/auth/password/change", r.authHandler.ChangePassword).Methods("POST")
	api.HandleFunc("/auth/logout", r.sessionHandler.Logout).Methods("POST")

	// Apply middleware stack
	middlewareStack := middleware.NewMiddlewareStack()
	handler := middlewareStack.ApplyToRouter(router)

	log.Println("✅ Routes and middleware configured successfully")
	return handler
}

// healthCheck handles GET /health
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	healthResp := response.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Data: map[string]interface{}{
			"uptime": "running",
		},
	}

	response.WriteJSON(w, http.StatusOK, healthResp)
}
