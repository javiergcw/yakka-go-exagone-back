package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	auth_rest "github.com/yakka-backend/internal/features/auth/delivery/rest"
	builder_rest "github.com/yakka-backend/internal/features/builder_profiles/delivery/rest"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	job_rest "github.com/yakka-backend/internal/features/jobs/delivery/rest"
	job_usecase "github.com/yakka-backend/internal/features/jobs/usecase"
	jobsite_rest "github.com/yakka-backend/internal/features/jobsites/delivery/rest"
	jobsite_db "github.com/yakka-backend/internal/features/jobsites/entity/database"
	labour_rest "github.com/yakka-backend/internal/features/labour_profiles/delivery/rest"
	experience_level_rest "github.com/yakka-backend/internal/features/masters/experience_levels/delivery/rest"
	job_requirement_rest "github.com/yakka-backend/internal/features/masters/job_requirements/delivery/rest"
	job_requirement_db "github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	job_type_rest "github.com/yakka-backend/internal/features/masters/job_types/delivery/rest"
	job_type_db "github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	license_rest "github.com/yakka-backend/internal/features/masters/licenses/delivery/rest"
	license_db "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	payment_constant_rest "github.com/yakka-backend/internal/features/masters/payment_constants/delivery/rest"
	payment_constant_usecase "github.com/yakka-backend/internal/features/masters/payment_constants/usecase"
	skill_category_rest "github.com/yakka-backend/internal/features/masters/skills/delivery/rest"
	skill_category_db "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	qualification_rest "github.com/yakka-backend/internal/features/qualifications/delivery/rest"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
)

// Router sets up the HTTP routes
type Router struct {
	authHandler                *auth_rest.AuthHandler
	sessionHandler             *auth_rest.SessionHandler
	passwordHandler            *auth_rest.PasswordHandler
	emailHandler               *auth_rest.EmailHandler
	labourProfileHandler       *labour_rest.LabourProfileHandler
	builderProfileHandler      *builder_rest.BuilderProfileHandler
	companyHandler             *builder_rest.CompanyHandler
	jobsiteHandler             *jobsite_rest.JobsiteHandler
	jobHandler                 *job_rest.JobHandler
	qualificationHandler       *qualification_rest.QualificationHandler
	labourQualificationHandler *qualification_rest.LabourQualificationHandler
	licenseHandler             *license_rest.LicenseHandler
	experienceLevelHandler     *experience_level_rest.ExperienceLevelHandler
	skillCategoryHandler       *skill_category_rest.SkillCategoryHandler
	skillSubcategoryHandler    *skill_category_rest.SkillSubcategoryHandler
	skillCompleteHandler       *skill_category_rest.SkillCompleteHandler
	jobRequirementHandler      *job_requirement_rest.JobRequirementHandler
	jobTypeHandler             *job_type_rest.JobTypeHandler
	paymentConstantHandler     *payment_constant_rest.PaymentConstantHandler
}

// NewRouter creates a new router
func NewRouter(
	authHandler *auth_rest.AuthHandler,
	sessionHandler *auth_rest.SessionHandler,
	passwordHandler *auth_rest.PasswordHandler,
	emailHandler *auth_rest.EmailHandler,
	labourProfileHandler *labour_rest.LabourProfileHandler,
	builderProfileHandler *builder_rest.BuilderProfileHandler,
	companyHandler *builder_rest.CompanyHandler,
	jobsiteHandler *jobsite_rest.JobsiteHandler,
	qualificationHandler *qualification_rest.QualificationHandler,
	labourQualificationHandler *qualification_rest.LabourQualificationHandler,
	jobUsecase job_usecase.JobUsecase,
	builderProfileRepo builder_db.BuilderProfileRepository,
	jobsiteRepo jobsite_db.JobsiteRepository,
	jobTypeRepo job_type_db.JobTypeRepository,
	licenseRepo license_db.LicenseRepository,
	paymentConstantUseCase payment_constant_usecase.PaymentConstantUsecase,
	jobRequirementRepo job_requirement_db.JobRequirementRepository,
	skillCategoryRepo skill_category_db.SkillCategoryRepository,
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository,
) *Router {
	return &Router{
		authHandler:                authHandler,
		sessionHandler:             sessionHandler,
		passwordHandler:            passwordHandler,
		emailHandler:               emailHandler,
		labourProfileHandler:       labourProfileHandler,
		builderProfileHandler:      builderProfileHandler,
		companyHandler:             companyHandler,
		jobsiteHandler:             jobsiteHandler,
		qualificationHandler:       qualificationHandler,
		labourQualificationHandler: labourQualificationHandler,
		jobHandler:                 job_rest.NewJobHandler(jobUsecase, builderProfileRepo, jobsiteRepo, jobTypeRepo, licenseRepo, jobRequirementRepo, skillCategoryRepo, skillSubcategoryRepo),
		licenseHandler:             license_rest.NewLicenseHandler(),
		experienceLevelHandler:     experience_level_rest.NewExperienceLevelHandler(),
		skillCategoryHandler:       skill_category_rest.NewSkillCategoryHandler(),
		skillSubcategoryHandler:    skill_category_rest.NewSkillSubcategoryHandler(),
		skillCompleteHandler:       skill_category_rest.NewSkillCompleteHandler(),
		jobRequirementHandler:      job_requirement_rest.NewJobRequirementHandler(jobRequirementRepo),
		jobTypeHandler:             job_type_rest.NewJobTypeHandler(jobTypeRepo),
		paymentConstantHandler:     payment_constant_rest.NewPaymentConstantHandler(paymentConstantUseCase),
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

	// Company endpoints (require license)
	api.Handle("/companies", middleware.LicenseMiddleware(http.HandlerFunc(r.companyHandler.CreateCompany))).Methods("POST")
	api.Handle("/companies", middleware.LicenseMiddleware(http.HandlerFunc(r.companyHandler.GetCompanies))).Methods("GET")

	// Master tables endpoints (require license)
	api.Handle("/licenses", middleware.LicenseMiddleware(http.HandlerFunc(r.licenseHandler.GetLicenses))).Methods("GET")
	api.Handle("/experience-levels", middleware.LicenseMiddleware(http.HandlerFunc(r.experienceLevelHandler.GetExperienceLevels))).Methods("GET")
	api.Handle("/skill-categories", middleware.LicenseMiddleware(http.HandlerFunc(r.skillCategoryHandler.GetSkillCategories))).Methods("GET")
	api.Handle("/skill-subcategories", middleware.LicenseMiddleware(http.HandlerFunc(r.skillSubcategoryHandler.GetSkillSubcategories))).Methods("GET")
	api.Handle("/skill-categories/{categoryId}/subcategories", middleware.LicenseMiddleware(http.HandlerFunc(r.skillSubcategoryHandler.GetSkillSubcategoriesByCategory))).Methods("GET")
	api.Handle("/skills", middleware.LicenseMiddleware(http.HandlerFunc(r.skillCompleteHandler.GetSkillsComplete))).Methods("GET")
	api.Handle("/qualifications", middleware.LicenseMiddleware(http.HandlerFunc(r.qualificationHandler.GetQualifications))).Methods("GET")
	api.Handle("/payment-constants", middleware.LicenseMiddleware(http.HandlerFunc(r.paymentConstantHandler.GetAllPaymentConstants))).Methods("GET")
	api.Handle("/job-requirements", middleware.LicenseMiddleware(http.HandlerFunc(r.jobRequirementHandler.GetJobRequirements))).Methods("GET")
	api.Handle("/job-types", middleware.LicenseMiddleware(http.HandlerFunc(r.jobTypeHandler.GetJobTypes))).Methods("GET")

	// Protected endpoints (require authentication only)
	api.Handle("/profiles/labour", middleware.AuthMiddleware(http.HandlerFunc(r.labourProfileHandler.CreateLabourProfile))).Methods("POST")
	api.Handle("/profiles/builder", middleware.AuthMiddleware(http.HandlerFunc(r.builderProfileHandler.CreateBuilderProfile))).Methods("POST")
	api.Handle("/auth/profile", middleware.AuthMiddleware(http.HandlerFunc(r.authHandler.GetProfile))).Methods("GET")

	// Builder endpoints (require builder role)
	api.Handle("/builder/companies", middleware.BuilderMiddleware(http.HandlerFunc(r.companyHandler.AssignCompany))).Methods("POST")
	api.Handle("/jobsites", middleware.BuilderMiddleware(http.HandlerFunc(r.jobsiteHandler.CreateJobsite))).Methods("POST")
	api.Handle("/jobsites", middleware.BuilderMiddleware(http.HandlerFunc(r.jobsiteHandler.GetJobsitesByBuilder))).Methods("GET")
	api.Handle("/jobsites/{id}", middleware.BuilderMiddleware(http.HandlerFunc(r.jobsiteHandler.GetJobsiteByID))).Methods("GET")
	api.Handle("/jobsites/{id}", middleware.BuilderMiddleware(http.HandlerFunc(r.jobsiteHandler.UpdateJobsite))).Methods("PUT")
	api.Handle("/jobsites/{id}", middleware.BuilderMiddleware(http.HandlerFunc(r.jobsiteHandler.DeleteJobsite))).Methods("DELETE")

	// Job endpoints (require builder role)
	api.Handle("/builder/jobs", middleware.BuilderMiddleware(http.HandlerFunc(r.jobHandler.CreateJob))).Methods("POST")
	api.Handle("/builder/jobs", middleware.BuilderMiddleware(http.HandlerFunc(r.jobHandler.GetMyJobs))).Methods("GET")
	api.Handle("/builder/jobs/{id}", middleware.BuilderMiddleware(http.HandlerFunc(r.jobHandler.GetBuilderJobDetail))).Methods("GET")
	api.Handle("/builder/jobs/{id}/visibility", middleware.BuilderMiddleware(http.HandlerFunc(r.jobHandler.UpdateJobVisibility))).Methods("PUT")
	api.Handle("/builder/applicants", middleware.BuilderMiddleware(http.HandlerFunc(r.jobHandler.GetBuilderApplicants))).Methods("GET")
	api.Handle("/builder/applicants", middleware.BuilderMiddleware(http.HandlerFunc(r.jobHandler.ProcessApplicantDecision))).Methods("POST")

	// Labour endpoints (require labour role)
	api.Handle("/labour/jobs", middleware.LabourMiddleware(http.HandlerFunc(r.jobHandler.GetLabourJobs))).Methods("GET")
	api.Handle("/labour/jobs/{id}", middleware.LabourMiddleware(http.HandlerFunc(r.jobHandler.GetLabourJobDetail))).Methods("GET")
	api.Handle("/labour/applicants", middleware.LabourMiddleware(http.HandlerFunc(r.jobHandler.GetLabourApplicants))).Methods("GET")
	api.Handle("/labour/applicants", middleware.LabourMiddleware(http.HandlerFunc(r.jobHandler.ApplyToJob))).Methods("POST")
	api.Handle("/labour/qualifications", middleware.LabourMiddleware(http.HandlerFunc(r.labourQualificationHandler.GetLabourQualifications))).Methods("GET")
	api.Handle("/labour/qualifications", middleware.LabourMiddleware(http.HandlerFunc(r.labourQualificationHandler.CreateLabourQualifications))).Methods("POST")
	api.Handle("/labour/qualifications", middleware.LabourMiddleware(http.HandlerFunc(r.labourQualificationHandler.UpdateLabourQualifications))).Methods("PUT")

	//labour endpoints

	/*
		no test
			api.HandleFunc("/auth/profile", r.authHandler.UpdateProfile).Methods("PUT")
			api.HandleFunc("/auth/password/change", r.authHandler.ChangePassword).Methods("POST")
			api.HandleFunc("/auth/logout", r.sessionHandler.Logout).Methods("POST")

			api.HandleFunc("/auth/refresh", r.sessionHandler.RefreshToken).Methods("POST")
		   	api.HandleFunc("/auth/password/reset", r.passwordHandler.RequestPasswordReset).Methods("POST")
		   	api.HandleFunc("/auth/password/reset/confirm", r.passwordHandler.ResetPassword).Methods("POST")
		   	api.HandleFunc("/auth/email/verify", r.emailHandler.VerifyEmail).Methods("POST")
	*/

	// Apply middleware stack (basic middleware only)
	handler := middlewareStack.ApplyToRouter(router)

	log.Println("✅ Routes and middleware configured successfully")
	return handler
}

// healthCheck handles GET /health
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	healthResp := response.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.3.4",
		Data: map[string]interface{}{
			"uptime":  "running",
			"license": "YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B", // License for master tables endpoints
		},
	}

	response.WriteJSON(w, http.StatusOK, healthResp)
}
