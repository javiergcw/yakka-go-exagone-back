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
	labour_rest "github.com/yakka-backend/internal/features/labour_profiles/delivery/rest"
	labour_db "github.com/yakka-backend/internal/features/labour_profiles/entity/database"
	labour_usecase "github.com/yakka-backend/internal/features/labour_profiles/usecase"
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

	// Initialize use cases
	authUserUseCase := auth_user_usecase.NewAuthUsecase(authUserRepo, builderRepo, labourRepo)
	authSessionUseCase := auth_session_usecase.NewSessionUsecase(authSessionRepo)
	authPasswordUseCase := auth_password_usecase.NewPasswordResetUsecase(authPasswordRepo)
	authEmailUseCase := auth_email_usecase.NewEmailVerificationUsecase(authEmailRepo, authUserRepo)
	labourSkillRepo := labour_db.NewLabourProfileSkillRepository(database.DB)
	userLicenseRepo := auth_user_db.NewUserLicenseRepository(database.DB)
	labourProfileUseCase := labour_usecase.NewLabourProfileUsecase(labourRepo, labourSkillRepo, userLicenseRepo, authUserRepo)
	builderProfileUseCase := builder_usecase.NewBuilderProfileUsecase(builderRepo, userLicenseRepo, authUserRepo)

	// Initialize handlers
	authHandler := auth_rest.NewAuthHandler(authUserUseCase, authEmailUseCase, builderProfileUseCase, labourProfileUseCase)
	sessionHandler := auth_rest.NewSessionHandler(authSessionUseCase)
	passwordHandler := auth_rest.NewPasswordHandler(authPasswordUseCase)
	emailHandler := auth_rest.NewEmailHandler(authEmailUseCase)
	labourProfileHandler := labour_rest.NewLabourProfileHandler(labourProfileUseCase)
	builderProfileHandler := builder_rest.NewBuilderProfileHandler(builderProfileUseCase)

	// Initialize router
	router := httpRouter.NewRouter(authHandler, sessionHandler, passwordHandler, emailHandler, labourProfileHandler, builderProfileHandler)
	httpRouter := router.SetupRoutes()

	// Start server
	fmt.Printf("🚀 Server starting on port %s\n", cfg.Server.Port)
	fmt.Printf("📋 Health check: http://localhost:%s/health\n", cfg.Server.Port)
	fmt.Printf("👥 Users API: http://localhost:%s/api/v1/users\n", cfg.Server.Port)
	fmt.Printf("🗄️  Database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, httpRouter))
}
