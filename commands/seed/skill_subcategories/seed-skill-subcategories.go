package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	skillRepo "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	skillSubcategoryModels "github.com/yakka-backend/internal/features/masters/skills/models"
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

	// Initialize repositories
	skillCategoryRepository := skillRepo.NewSkillCategoryRepository(database.DB)
	skillSubcategoryRepository := skillRepo.NewSkillSubcategoryRepository(database.DB)

	// Get all categories to map names to IDs
	categories, err := skillCategoryRepository.GetAll(context.TODO())
	if err != nil {
		log.Fatalf("Failed to get categories: %v", err)
	}

	// Create category name to ID mapping
	categoryMap := make(map[string]uuid.UUID)
	for _, category := range categories {
		categoryMap[category.Name] = category.ID
	}

	// Sample skill subcategories data
	subcategoriesData := map[string][]string{
		"Coach": {
			"Soccer", "Rugby Union", "Rugby League", "Australian Rules Football (AFL)",
			"Basketball", "Netball", "Tennis", "Swimming", "Athletics", "Cricket",
			"Hockey", "Volleyball", "Gymnastics", "Baseball", "Softball", "Surfing",
			"Cycling", "Boxing", "Martial Arts", "Golf", "Other",
		},
		"Referee / Umpire / Match Official": {
			"Soccer", "Rugby Union", "Rugby League", "AFL", "Basketball", "Netball",
			"Tennis", "Volleyball", "Handball", "Cricket", "Baseball", "Softball",
			"Hockey", "Water Polo", "Boxing", "Martial Arts", "Athletics", "Other",
		},
		"Personal Trainer (PT)": {
			"Fitness", "Gym", "CrossFit", "Functional Training", "Athletics", "Swimming",
			"Cycling", "Boxing", "Martial Arts", "Yoga", "Pilates", "Other",
		},
		"Strength & Conditioning Coach": {
			"Soccer", "Rugby Union", "Rugby League", "AFL", "Basketball", "Netball",
			"Tennis", "Swimming", "Athletics", "Cricket", "Cycling", "Combat Sports",
			"Gymnastics", "Other",
		},
		"Fitness Instructor / Group Trainer": {
			"Gym", "CrossFit", "Pilates", "Yoga", "Dance Sport", "Aquatic Fitness",
			"Martial Arts", "Other",
		},
		"High Performance Manager": {
			"Soccer", "Rugby", "AFL", "Basketball", "Netball", "Swimming", "Athletics",
			"Cricket", "Cycling", "Combat Sports", "Other",
		},
		"Athlete Development Officer": {
			"Soccer", "Rugby", "AFL", "Basketball", "Netball", "Tennis", "Swimming",
			"Athletics", "Gymnastics", "Other",
		},
	}

	// Insert subcategories
	for categoryName, sports := range subcategoriesData {
		categoryID, exists := categoryMap[categoryName]
		if !exists {
			log.Printf("Category %s not found, skipping...", categoryName)
			continue
		}

		for _, sport := range sports {
			subcategory := &skillSubcategoryModels.SkillSubcategory{
				CategoryID:  categoryID,
				Name:        sport,
				Description: "Sport specialization for " + categoryName,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if err := skillSubcategoryRepository.Create(context.TODO(), subcategory); err != nil {
				log.Printf("Failed to create subcategory %s for category %s: %v", sport, categoryName, err)
			} else {
				log.Printf("✅ Created subcategory: %s (%s)", sport, categoryName)
			}
		}
	}

	log.Println("✅ Skill subcategory seeding completed successfully!")
}
