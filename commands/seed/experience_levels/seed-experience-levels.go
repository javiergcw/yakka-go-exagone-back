package main

import (
	"context"
	"log"
	"time"

	experienceLevelRepo "github.com/yakka-backend/internal/features/masters/experience_levels/entity/database"
	experienceLevelModels "github.com/yakka-backend/internal/features/masters/experience_levels/models"
	"github.com/yakka-backend/internal/infrastructure/config"
	"github.com/yakka-backend/internal/infrastructure/database"
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

	// Initialize experience level repository
	experienceLevelRepository := experienceLevelRepo.NewExperienceLevelRepository(database.DB)

	// Sample experience levels
	experienceLevels := []*experienceLevelModels.ExperienceLevel{
		{
			Name:        "Less than 6 months",
			Description: "Experience level for workers with less than 6 months of experience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "6-12 months",
			Description: "Experience level for workers with 6 to 12 months of experience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "1-2 years",
			Description: "Experience level for workers with 1 to 2 years of experience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "2-5 years",
			Description: "Experience level for workers with 2 to 5 years of experience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "More than 5 years",
			Description: "Experience level for workers with more than 5 years of experience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert experience levels
	for _, experienceLevel := range experienceLevels {
		if err := experienceLevelRepository.Create(context.TODO(), experienceLevel); err != nil {
			log.Printf("Failed to create experience level %s: %v", experienceLevel.Name, err)
		} else {
			log.Printf("✅ Created experience level: %s", experienceLevel.Name)
		}
	}

	log.Println("✅ Experience level seeding completed successfully!")
}
