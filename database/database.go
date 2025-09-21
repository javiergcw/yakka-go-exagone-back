package database

import (
	"fmt"
	"log"

	"github.com/yakka-backend/config"
	auth_email_models "github.com/yakka-backend/internal/features/auth/email_verification/models"
	auth_password_models "github.com/yakka-backend/internal/features/auth/password_reset/models"
	auth_user_models "github.com/yakka-backend/internal/features/auth/user/models"
	auth_session_models "github.com/yakka-backend/internal/features/auth/user_session/models"
	user_models "github.com/yakka-backend/internal/features/users/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	err := DB.AutoMigrate(
		&user_models.User{},
		&auth_user_models.User{},
		&auth_session_models.Session{},
		&auth_password_models.PasswordReset{},
		&auth_email_models.EmailVerification{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}
