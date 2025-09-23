package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	auth_rest "github.com/yakka-backend/internal/features/auth/delivery/rest"
	builder_rest "github.com/yakka-backend/internal/features/builder_profiles/delivery/rest"
	labour_rest "github.com/yakka-backend/internal/features/labour_profiles/delivery/rest"
	experience_level_rest "github.com/yakka-backend/internal/features/masters/experience_levels/delivery/rest"
	license_rest "github.com/yakka-backend/internal/features/masters/licenses/delivery/rest"
	skill_category_rest "github.com/yakka-backend/internal/features/masters/skills/delivery/rest"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
)

// Router sets up the HTTP routes
type Router struct {
	authHandler             *auth_rest.AuthHandler
	sessionHandler          *auth_rest.SessionHandler
	passwordHandler         *auth_rest.PasswordHandler
	emailHandler            *auth_rest.EmailHandler
	labourProfileHandler    *labour_rest.LabourProfileHandler
	builderProfileHandler   *builder_rest.BuilderProfileHandler
	licenseHandler          *license_rest.LicenseHandler
	experienceLevelHandler  *experience_level_rest.ExperienceLevelHandler
	skillCategoryHandler    *skill_category_rest.SkillCategoryHandler
	skillSubcategoryHandler *skill_category_rest.SkillSubcategoryHandler
	skillCompleteHandler    *skill_category_rest.SkillCompleteHandler
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
		authHandler:             authHandler,
		sessionHandler:          sessionHandler,
		passwordHandler:         passwordHandler,
		emailHandler:            emailHandler,
		labourProfileHandler:    labourProfileHandler,
		builderProfileHandler:   builderProfileHandler,
		licenseHandler:          license_rest.NewLicenseHandler(),
		experienceLevelHandler:  experience_level_rest.NewExperienceLevelHandler(),
		skillCategoryHandler:    skill_category_rest.NewSkillCategoryHandler(),
		skillSubcategoryHandler: skill_category_rest.NewSkillSubcategoryHandler(),
		skillCompleteHandler:    skill_category_rest.NewSkillCompleteHandler(),
	}
}

// SetupRoutes configures all the routes
func (r *Router) SetupRoutes() http.Handler {
	router := mux.NewRouter()
	middlewareStack := middleware.NewMiddlewareStack()

	// Health check endpoint (public)
	router.HandleFunc("/health", r.healthCheck).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public auth endpoints (no middleware)
	api.HandleFunc("/auth/register", r.authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", r.authHandler.Login).Methods("POST")
	/*
		no test
		api.HandleFunc("/auth/refresh", r.sessionHandler.RefreshToken).Methods("POST")
		   	api.HandleFunc("/auth/password/reset", r.passwordHandler.RequestPasswordReset).Methods("POST")
		   	api.HandleFunc("/auth/password/reset/confirm", r.passwordHandler.ResetPassword).Methods("POST")
		   	api.HandleFunc("/auth/email/verify", r.emailHandler.VerifyEmail).Methods("POST") */

	// Master tables endpoints (require license only)
	licenseAPI := api.PathPrefix("").Subrouter()
	licenseAPI.Handle("/licenses", http.HandlerFunc(r.licenseHandler.GetLicenses)).Methods("GET")
	licenseAPI.Handle("/experience-levels", http.HandlerFunc(r.experienceLevelHandler.GetExperienceLevels)).Methods("GET")
	licenseAPI.Handle("/skill-categories", http.HandlerFunc(r.skillCategoryHandler.GetSkillCategories)).Methods("GET")
	licenseAPI.Handle("/skill-subcategories", http.HandlerFunc(r.skillSubcategoryHandler.GetSkillSubcategories)).Methods("GET")
	licenseAPI.Handle("/skill-categories/{categoryId}/subcategories", http.HandlerFunc(r.skillSubcategoryHandler.GetSkillSubcategoriesByCategory)).Methods("GET")
	licenseAPI.Handle("/skills", http.HandlerFunc(r.skillCompleteHandler.GetSkillsComplete)).Methods("GET")

	// Protected endpoints (require authentication)
	authAPI := api.PathPrefix("").Subrouter()
	authAPI.HandleFunc("/auth/profile", r.authHandler.GetProfile).Methods("GET")
	/*
		no test
			authAPI.HandleFunc("/auth/profile", r.authHandler.UpdateProfile).Methods("PUT")
			authAPI.HandleFunc("/auth/password/change", r.authHandler.ChangePassword).Methods("POST")
			authAPI.HandleFunc("/auth/logout", r.sessionHandler.Logout).Methods("POST")
			authAPI.HandleFunc("/profiles/labour", r.labourProfileHandler.CreateLabourProfile).Methods("POST")
			authAPI.HandleFunc("/profiles/builder", r.builderProfileHandler.CreateBuilderProfile).Methods("POST") */

	// Apply middleware to specific subrouters
	licenseHandler := middlewareStack.ApplyWithLicense(licenseAPI)
	authHandler := middlewareStack.ApplyWithAuth(authAPI)

	// Mount the subrouters with their respective middleware
	api.PathPrefix("").Handler(licenseHandler)
	api.PathPrefix("").Handler(authHandler)

	// Apply basic middleware to the main router
	handler := middlewareStack.ApplyPublic(router)

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
			"uptime":  "running",
			"license": "YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B", // License for master tables endpoints
		},
	}

	response.WriteJSON(w, http.StatusOK, healthResp)
}
