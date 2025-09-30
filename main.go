package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	auth_rest "github.com/yakka-backend/internal/features/auth/delivery/rest"
	auth_email_db "github.com/yakka-backend/internal/features/auth/email_verification/entity/database"
	auth_email_usecase "github.com/yakka-backend/internal/features/auth/email_verification/usecase"
	auth_password_db "github.com/yakka-backend/internal/features/auth/password_reset/entity/database"
	auth_password_usecase "github.com/yakka-backend/internal/features/auth/password_reset/usecase"
	auth_user_db "github.com/yakka-backend/internal/features/auth/user/entity/database"
	auth_user_usecase "github.com/yakka-backend/internal/features/auth/user/usecase"
	auth_session_db "github.com/yakka-backend/internal/features/auth/user_session/entity/database"
	auth_session_usecase "github.com/yakka-backend/internal/features/auth/user_session/usecase"
	builder_rest "github.com/yakka-backend/internal/features/builder_profiles/delivery/rest"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	builder_usecase "github.com/yakka-backend/internal/features/builder_profiles/usecase"
	job_application_db "github.com/yakka-backend/internal/features/job_applications/entity/database"
	job_assignment_db "github.com/yakka-backend/internal/features/job_assignments/entity/database"
	job_db "github.com/yakka-backend/internal/features/jobs/entity/database"
	job_usecase "github.com/yakka-backend/internal/features/jobs/usecase"
	jobsite_rest "github.com/yakka-backend/internal/features/jobsites/delivery/rest"
	jobsite_db "github.com/yakka-backend/internal/features/jobsites/entity/database"
	jobsite_usecase "github.com/yakka-backend/internal/features/jobsites/usecase"
	labour_rest "github.com/yakka-backend/internal/features/labour_profiles/delivery/rest"
	labour_db "github.com/yakka-backend/internal/features/labour_profiles/entity/database"
	labour_usecase "github.com/yakka-backend/internal/features/labour_profiles/usecase"
	experience_db "github.com/yakka-backend/internal/features/masters/experience_levels/entity/database"
	job_requirement_db "github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	job_type_db "github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	license_db "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	payment_constant_db "github.com/yakka-backend/internal/features/masters/payment_constants/entity/database"
	payment_constant_usecase "github.com/yakka-backend/internal/features/masters/payment_constants/usecase"
	skill_db "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	"github.com/yakka-backend/internal/infrastructure/config"
	"github.com/yakka-backend/internal/infrastructure/database"
	httpRouter "github.com/yakka-backend/internal/infrastructure/http"
)

func main() {
	// Parse command line flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// If migrate flag is set, run migrations and exit
	if *migrateFlag {
		log.Println("🚀 Running database migrations...")
		if err := database.Migrate(); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("✅ Database migration completed successfully!")
		return
	}

	// Initialize repositories
	authUserRepo := auth_user_db.NewUserRepository(database.DB)
	authSessionRepo := auth_session_db.NewSessionRepository(database.DB)
	authPasswordRepo := auth_password_db.NewPasswordResetRepository(database.DB)
	authEmailRepo := auth_email_db.NewEmailVerificationRepository(database.DB)
	builderRepo := builder_db.NewBuilderProfileRepository(database.DB)
	labourRepo := labour_db.NewLabourProfileRepository(database.DB)
	jobsiteRepo := jobsite_db.NewJobsiteRepositoryImpl(database.DB)

	// Initialize use cases
	authUserUseCase := auth_user_usecase.NewAuthUsecase(authUserRepo, builderRepo, labourRepo)
	authSessionUseCase := auth_session_usecase.NewSessionUsecase(authSessionRepo)
	authPasswordUseCase := auth_password_usecase.NewPasswordResetUsecase(authPasswordRepo)
	authEmailUseCase := auth_email_usecase.NewEmailVerificationUsecase(authEmailRepo, authUserRepo)
	labourSkillRepo := labour_db.NewLabourProfileSkillRepository(database.DB)
	userLicenseRepo := auth_user_db.NewUserLicenseRepository(database.DB)
	licenseRepo := license_db.NewLicenseRepository(database.DB)
	skillCategoryRepo := skill_db.NewSkillCategoryRepository(database.DB)
	skillSubcategoryRepo := skill_db.NewSkillSubcategoryRepository(database.DB)
	experienceRepo := experience_db.NewExperienceLevelRepository(database.DB)
	jobRequirementRepo := job_requirement_db.NewJobRequirementRepository(database.DB)
	jobTypeRepo := job_type_db.NewJobTypeRepository(database.DB)
	paymentConstantRepo := payment_constant_db.NewPaymentConstantRepository(database.DB)

	// Job repositories
	jobRepo := job_db.NewJobRepository(database.DB)
	jobLicenseRepo := job_db.NewJobLicenseRepository(database.DB)
	jobSkillRepo := job_db.NewJobSkillRepository(database.DB)
	jobJobRequirementRepo := job_db.NewJobJobRequirementRepository(database.DB)

	// Job Application repositories
	jobApplicationRepo := job_application_db.NewJobApplicationRepository(database.DB)

	// Job Assignment repositories
	jobAssignmentRepo := job_assignment_db.NewJobAssignmentRepository(database.DB)

	labourProfileUseCase := labour_usecase.NewLabourProfileUsecase(labourRepo, labourSkillRepo, userLicenseRepo, authUserRepo, licenseRepo, skillCategoryRepo, skillSubcategoryRepo, experienceRepo)
	builderProfileUseCase := builder_usecase.NewBuilderProfileUsecase(builderRepo, userLicenseRepo, authUserRepo, licenseRepo)
	jobsiteUseCase := jobsite_usecase.NewJobsiteUsecaseImpl(jobsiteRepo)
	paymentConstantUseCase := payment_constant_usecase.NewPaymentConstantUsecase(paymentConstantRepo)
	// jobApplicationUseCase := job_application_usecase.NewJobApplicationUsecase(jobApplicationRepo) // Available for future use
	// jobAssignmentUseCase := job_assignment_usecase.NewJobAssignmentUsecase(jobAssignmentRepo) // Available for future use
	jobUseCase := job_usecase.NewJobUsecase(jobRepo, jobLicenseRepo, jobSkillRepo, jobJobRequirementRepo, jobRequirementRepo, builderRepo, jobsiteRepo, jobTypeRepo, jobApplicationRepo, jobAssignmentRepo, licenseRepo, skillCategoryRepo, skillSubcategoryRepo)

	// Initialize handlers
	authHandler := auth_rest.NewAuthHandler(authUserUseCase, authEmailUseCase, builderProfileUseCase, labourProfileUseCase)
	sessionHandler := auth_rest.NewSessionHandler(authSessionUseCase)
	passwordHandler := auth_rest.NewPasswordHandler(authPasswordUseCase)
	emailHandler := auth_rest.NewEmailHandler(authEmailUseCase)
	labourProfileHandler := labour_rest.NewLabourProfileHandler(labourProfileUseCase)
	builderProfileHandler := builder_rest.NewBuilderProfileHandler(builderProfileUseCase)
	jobsiteHandler := jobsite_rest.NewJobsiteHandler(jobsiteUseCase)
	// jobApplicationHandler := job_application_rest.NewJobApplicationHandler(jobApplicationUseCase) // Available for future use
	// jobAssignmentHandler := job_assignment_rest.NewJobAssignmentHandler(jobAssignmentUseCase) // Available for future use

	// Initialize router
	router := httpRouter.NewRouter(authHandler, sessionHandler, passwordHandler, emailHandler, labourProfileHandler, builderProfileHandler, jobsiteHandler, jobUseCase, builderRepo, jobsiteRepo, jobTypeRepo, licenseRepo, paymentConstantUseCase, jobRequirementRepo)
	httpRouter := router.SetupRoutes()

	// Start server
	fmt.Printf("🚀 Server starting on port %s\n", cfg.Server.Port)
	fmt.Printf("📋 Health check: http://localhost:%s/health\n", cfg.Server.Port)
	fmt.Printf("👥 Users API: http://localhost:%s/api/v1/users\n", cfg.Server.Port)
	fmt.Printf("🗄️  Database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	// NEW: usa PORT y escucha en 0.0.0.0
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: httpRouter,
	}

	log.Fatal(srv.ListenAndServe())
}
