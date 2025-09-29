package main

import (
	"context"
	"log"
	"time"

	licenseRepo "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	licenseModels "github.com/yakka-backend/internal/features/masters/licenses/models"
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

	// Initialize license repository
	licenseRepository := licenseRepo.NewLicenseRepository(database.DB)

	// Sample licenses
	licenses := []*licenseModels.License{
		{
			Name:        "Passport",
			Description: "Official passport for international identification and travel",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Photo ID",
			Description: "Government-issued photo identification document",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Driver License",
			Description: "Official license to operate motor vehicles legally",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "WWCC",
			Description: "Working with Children Check - required for working with children",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Visa",
			Description: "Valid visa for work authorization in Australia",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Police Check",
			Description: "National Police Certificate for background verification",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "First Aid Training",
			Description: "Current first aid certification for workplace safety",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "White Card",
			Description: "Construction Industry Safety Card - required for construction work",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "RSA",
			Description: "Responsible Service of Alcohol certification",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert licenses
	for _, license := range licenses {
		if err := licenseRepository.Create(context.TODO(), license); err != nil {
			log.Printf("Failed to create license %s: %v", license.Name, err)
		} else {
			log.Printf("✅ Created license: %s", license.Name)
		}
	}

	log.Println("✅ License seeding completed successfully!")
}
