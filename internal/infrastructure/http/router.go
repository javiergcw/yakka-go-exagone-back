package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	auth_rest "github.com/yakka-backend/internal/features/auth/delivery/rest"
	builder_rest "github.com/yakka-backend/internal/features/builder_profiles/delivery/rest"
	labour_rest "github.com/yakka-backend/internal/features/labour_profiles/delivery/rest"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
)

// Router sets up the HTTP routes
type Router struct {
	authHandler         *auth_rest.AuthHandler
	sessionHandler      *auth_rest.SessionHandler
	passwordHandler     *auth_rest.PasswordHandler
	emailHandler        *auth_rest.EmailHandler
	labourProfileHandler *labour_rest.LabourProfileHandler
	builderProfileHandler *builder_rest.BuilderProfileHandler
}

// NewRouter creates a new router
func NewRouter(
	authHandler *auth_rest.AuthHandler,
	sessionHandler *auth_rest.SessionHandler,
	passwordHandler *auth_rest.PasswordHandler,
	emailHandler *auth_rest.EmailHandler,
	labourProfileHandler *labour_rest.LabourProfileHandler,
	builderProfileHandler *builder_rest.BuilderProfileHandler,
) *Router {
	return &Router{
		authHandler:          authHandler,
		sessionHandler:       sessionHandler,
		passwordHandler:      passwordHandler,
		emailHandler:         emailHandler,
		labourProfileHandler: labourProfileHandler,
		builderProfileHandler: builderProfileHandler,
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

	// Profile endpoints (protected)
	api.HandleFunc("/profiles/labour", r.labourProfileHandler.CreateLabourProfile).Methods("POST")
	api.HandleFunc("/profiles/builder", r.builderProfileHandler.CreateBuilderProfile).Methods("POST")

	// Master tables endpoints (require license)
	api.Handle("/licenses", middleware.LicenseMiddleware(http.HandlerFunc(r.getLicenses))).Methods("GET")
	api.Handle("/skills", middleware.LicenseMiddleware(http.HandlerFunc(r.getSkills))).Methods("GET")

	// Apply middleware stack (includes auth middleware)
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
			"license": "YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B", // License for master tables endpoints
		},
	}

	response.WriteJSON(w, http.StatusOK, healthResp)
}

// getLicenses returns all licenses (requires license header)
func (r *Router) getLicenses(w http.ResponseWriter, req *http.Request) {
	// This would typically fetch from database
	licenses := []map[string]interface{}{
		{
			"id":          "1",
			"name":        "Licencia de Conducir",
			"description": "Permiso para conducir vehículos automotores",
		},
		{
			"id":          "2", 
			"name":        "Licencia de Construcción",
			"description": "Permiso para realizar trabajos de construcción",
		},
		{
			"id":          "3",
			"name":        "Licencia de Electricista", 
			"description": "Certificación para trabajos eléctricos",
		},
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    licenses,
		"message": "Licenses retrieved successfully",
	})
}

// getSkills returns all skills (requires license header)
func (r *Router) getSkills(w http.ResponseWriter, req *http.Request) {
	// This would typically fetch from database
	skills := []map[string]interface{}{
		{
			"id":          "1",
			"name":        "Albañilería",
			"description": "Construcción con ladrillos, bloques y mortero",
		},
		{
			"id":          "2",
			"name":        "Carpintería", 
			"description": "Trabajos con madera y estructuras de madera",
		},
		{
			"id":          "3",
			"name":        "Electricidad",
			"description": "Instalaciones y reparaciones eléctricas",
		},
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    skills,
		"message": "Skills retrieved successfully",
	})
}
