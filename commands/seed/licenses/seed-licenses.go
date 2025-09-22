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
			Name:        "Driving Licence",
			Description: "Official license to operate motor vehicles legally",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Work with Children",
			Description: "Certification required to work with children in various settings",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS20122 – Certificate II in Sport and Recreation",
			Description: "Entry-level qualification for sport and recreation industry",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS20321 – Certificate II in Sport Coaching",
			Description: "Foundation qualification for sport coaching roles",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS30122 – Certificate III in Sport, Aquatics and Recreation",
			Description: "Intermediate qualification in sport, aquatics and recreation management",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS30521 – Certificate III in Sport Coaching",
			Description: "Advanced coaching qualification for sport professionals",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS40321 – Certificate IV in Sport Coaching",
			Description: "Senior coaching qualification for leadership roles",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS50321 – Diploma of Sport",
			Description: "Comprehensive diploma program in sport management and coaching",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "SIS50122 – Diploma of Sport, Aquatics and Recreation Management",
			Description: "Specialized diploma in sport, aquatics and recreation management",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Diploma of Sport & Exercise Science (provider specific)",
			Description: "Provider-specific diploma in sport and exercise science",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Bachelor Degree in Sport Science",
			Description: "Undergraduate degree in sport science and related fields",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Bachelor Honours",
			Description: "Honours degree in sport science with research component",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Master's Degree in Sport",
			Description: "Postgraduate degree in sport science and management",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Doctoral Degree (PhD) in Sport",
			Description: "Highest academic qualification in sport science and research",
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
