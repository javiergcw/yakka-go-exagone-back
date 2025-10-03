package database

import (
	"fmt"
	"log"

	emailVerificationModels "github.com/yakka-backend/internal/features/auth/email_verification/models"
	passwordResetModels "github.com/yakka-backend/internal/features/auth/password_reset/models"
	authUserModels "github.com/yakka-backend/internal/features/auth/user/models"
	userSessionModels "github.com/yakka-backend/internal/features/auth/user_session/models"
	builderProfileModels "github.com/yakka-backend/internal/features/builder_profiles/models"
	jobApplicationModels "github.com/yakka-backend/internal/features/job_applications/models"
	jobAssignmentModels "github.com/yakka-backend/internal/features/job_assignments/models"
	jobModels "github.com/yakka-backend/internal/features/jobs/models"
	jobsiteModels "github.com/yakka-backend/internal/features/jobsites/models"
	labourProfileModels "github.com/yakka-backend/internal/features/labour_profiles/models"
	experienceLevelModels "github.com/yakka-backend/internal/features/masters/experience_levels/models"
	jobRequirementModels "github.com/yakka-backend/internal/features/masters/job_requirements/models"
	jobTypeModels "github.com/yakka-backend/internal/features/masters/job_types/models"
	licenseModels "github.com/yakka-backend/internal/features/masters/licenses/models"
	paymentConstantModels "github.com/yakka-backend/internal/features/masters/payment_constants/models"
	skillModels "github.com/yakka-backend/internal/features/masters/skills/models"
	qualificationModels "github.com/yakka-backend/internal/features/qualifications/models"
	"github.com/yakka-backend/internal/infrastructure/config"
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
		Logger: logger.Default.LogMode(logger.Warn), // Reduce logging for performance
		// Performance optimizations
		PrepareStmt:                              true, // Enable prepared statements
		DisableForeignKeyConstraintWhenMigrating: true, // Faster migrations
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool for better performance
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)   // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)  // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(0) // Connection lifetime (0 = unlimited)

	log.Println("✅ Database connected successfully")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Create custom types first
	err := createCustomTypes()
	if err != nil {
		return fmt.Errorf("failed to create custom types: %w", err)
	}

	// Auto-migrate all models
	err = DB.AutoMigrate(
		// Core user models
		&authUserModels.User{},
		&userSessionModels.Session{},

		// Authentication models
		&emailVerificationModels.EmailVerification{},
		&passwordResetModels.PasswordReset{},

		// Profile models
		&builderProfileModels.BuilderProfile{},
		&builderProfileModels.Company{},
		&labourProfileModels.LabourProfile{},
		&labourProfileModels.LabourProfileSkill{},

		// Jobsite models
		&jobsiteModels.Jobsite{},

		// User license models
		&authUserModels.UserLicense{},

		// License models
		&licenseModels.License{},

		// Skill category models
		&skillModels.SkillCategory{},

		// Skill subcategory models
		&skillModels.SkillSubcategory{},

		// Experience level models
		&experienceLevelModels.ExperienceLevel{},

		// Job requirement models
		&jobRequirementModels.JobRequirement{},

		// Job type models
		&jobTypeModels.JobType{},

		// Payment constant models
		&paymentConstantModels.PaymentConstant{},

		// Job models
		&jobModels.Job{},
		&jobModels.JobLicense{},
		&jobModels.JobSkill{},
		&jobModels.JobJobRequirement{},

		// Job Application models
		&jobApplicationModels.JobApplication{},

		// Job Assignment models
		&jobAssignmentModels.JobAssignment{},

		// Qualification models
		&qualificationModels.SportsQualification{},
		&qualificationModels.Qualification{},
		&qualificationModels.LabourProfileQualification{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// createCustomTypes creates PostgreSQL custom types
func createCustomTypes() error {
	// Create or update user_status enum
	err := DB.Exec(`
		DO $$ BEGIN
			-- Try to create the enum first
			CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending', 'banned');
		EXCEPTION
			WHEN duplicate_object THEN 
				-- If enum exists, add new values if they don't exist
				BEGIN
					ALTER TYPE user_status ADD VALUE IF NOT EXISTS 'banned';
				EXCEPTION
					WHEN duplicate_object THEN null;
				END;
		END $$;
	`).Error
	if err != nil {
		return err
	}

	// Create or update user_role enum
	err = DB.Exec(`
		DO $$ BEGIN
			-- Try to create the enum first
			CREATE TYPE user_role AS ENUM ('user', 'admin', 'builder', 'labour');
		EXCEPTION
			WHEN duplicate_object THEN 
				-- If enum exists, add new values if they don't exist
				BEGIN
					ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'builder';
					ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'labour';
				EXCEPTION
					WHEN duplicate_object THEN null;
				END;
		END $$;
	`).Error
	if err != nil {
		return err
	}

	log.Println("✅ Custom types created successfully")
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
