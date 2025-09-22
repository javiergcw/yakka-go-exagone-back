package main

import (
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
			Name:        "Licencia de Conducir",
			Description: "Permiso para conducir vehículos automotores",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Licencia de Construcción",
			Description: "Permiso para realizar trabajos de construcción",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Licencia de Electricista",
			Description: "Certificación para trabajos eléctricos",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Licencia de Plomería",
			Description: "Certificación para trabajos de plomería",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Licencia de Seguridad",
			Description: "Certificación en seguridad laboral",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert licenses
	for _, license := range licenses {
		if err := licenseRepository.Create(nil, license); err != nil {
			log.Printf("Failed to create license %s: %v", license.Name, err)
		} else {
			log.Printf("✅ Created license: %s", license.Name)
		}
	}

	log.Println("✅ License seeding completed successfully!")
}
