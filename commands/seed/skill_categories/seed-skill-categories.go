package main

import (
	"context"
	"log"
	"time"

	skillCategoryRepo "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	skillCategoryModels "github.com/yakka-backend/internal/features/masters/skills/models"
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

	// Initialize skill category repository
	skillCategoryRepository := skillCategoryRepo.NewSkillCategoryRepository(database.DB)

	// Sample skill categories
	categories := []*skillCategoryModels.SkillCategory{
		{
			Name:        "Coach",
			Description: "Professional coaching roles in various sports",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Referee / Umpire / Match Official",
			Description: "Officiating roles in various sports",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Personal Trainer (PT)",
			Description: "Personal training and fitness coaching",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Strength & Conditioning Coach",
			Description: "Specialized strength and conditioning coaching",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Fitness Instructor / Group Trainer",
			Description: "Group fitness and training instruction",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "High Performance Manager",
			Description: "High performance sports management",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Athlete Development Officer",
			Description: "Athlete development and talent identification",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert skill categories
	for _, category := range categories {
		if err := skillCategoryRepository.Create(context.TODO(), category); err != nil {
			log.Printf("Failed to create skill category %s: %v", category.Name, err)
		} else {
			log.Printf("✅ Created skill category: %s", category.Name)
		}
	}

	log.Println("✅ Skill category seeding completed successfully!")
}
