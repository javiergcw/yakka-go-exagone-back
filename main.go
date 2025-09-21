package main

import (
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
	"github.com/yakka-backend/internal/infrastructure/config"
	"github.com/yakka-backend/internal/infrastructure/database"
	httpRouter "github.com/yakka-backend/internal/infrastructure/http"
)

func main() {
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

	// Initialize repositories
	authUserRepo := auth_user_db.NewUserRepository(database.DB)
	authSessionRepo := auth_session_db.NewSessionRepository(database.DB)
	authPasswordRepo := auth_password_db.NewPasswordResetRepository(database.DB)
	authEmailRepo := auth_email_db.NewEmailVerificationRepository(database.DB)

	// Initialize use cases
	authUserUseCase := auth_user_usecase.NewAuthUsecase(authUserRepo)
	authSessionUseCase := auth_session_usecase.NewSessionUsecase(authSessionRepo)
	authPasswordUseCase := auth_password_usecase.NewPasswordResetUsecase(authPasswordRepo)
	authEmailUseCase := auth_email_usecase.NewEmailVerificationUsecase(authEmailRepo, authUserRepo)

	// Initialize handlers
	authHandler := auth_rest.NewAuthHandler(authUserUseCase, authEmailUseCase)
	sessionHandler := auth_rest.NewSessionHandler(authSessionUseCase)
	passwordHandler := auth_rest.NewPasswordHandler(authPasswordUseCase)
	emailHandler := auth_rest.NewEmailHandler(authEmailUseCase)

	// Initialize router
	router := httpRouter.NewRouter(authHandler, sessionHandler, passwordHandler, emailHandler)
	httpRouter := router.SetupRoutes()

	// Start server
	fmt.Printf("🚀 Server starting on port %s\n", cfg.Server.Port)
	fmt.Printf("📋 Health check: http://localhost:%s/health\n", cfg.Server.Port)
	fmt.Printf("👥 Users API: http://localhost:%s/api/v1/users\n", cfg.Server.Port)
	fmt.Printf("🗄️  Database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, httpRouter))
}
