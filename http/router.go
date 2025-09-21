package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yakka-backend/http/middleware"
	authRest "github.com/yakka-backend/internal/features/auth/delivery/rest"
	"github.com/yakka-backend/internal/shared/response"
)

// Router sets up the HTTP routes
type Router struct {
	authHandler     *authRest.AuthHandler
	sessionHandler  *authRest.SessionHandler
	emailHandler    *authRest.EmailHandler
	passwordHandler *authRest.PasswordHandler
}

// NewRouter creates a new router
func NewRouter(
	authHandler *authRest.AuthHandler,
	sessionHandler *authRest.SessionHandler,
	emailHandler *authRest.EmailHandler,
	passwordHandler *authRest.PasswordHandler,
) *Router {
	return &Router{
		authHandler:     authHandler,
		sessionHandler:  sessionHandler,
		emailHandler:    emailHandler,
		passwordHandler: passwordHandler,
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
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", r.authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", r.authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", r.sessionHandler.RefreshToken).Methods("POST")
	auth.HandleFunc("/logout", r.sessionHandler.Logout).Methods("POST")

	// Email verification endpoints (public)
	auth.HandleFunc("/email/verify", r.emailHandler.VerifyEmail).Methods("POST")

	// Password reset endpoints (public)
	auth.HandleFunc("/password/request-reset", r.passwordHandler.RequestPasswordReset).Methods("POST")
	auth.HandleFunc("/password/reset", r.passwordHandler.ResetPassword).Methods("POST")

	// Protected endpoints (require authentication)
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// User profile endpoints (protected)
	protected.HandleFunc("/profile", r.authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", r.authHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/change-password", r.authHandler.ChangePassword).Methods("POST")

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
