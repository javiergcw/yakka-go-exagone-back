package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	"github.com/yakka-backend/internal/features/masters/job_requirements/models"
	"github.com/yakka-backend/internal/infrastructure/config"
	infraDB "github.com/yakka-backend/internal/infrastructure/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	err = infraDB.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get database instance
	db := infraDB.DB

	// Create repository
	jobRequirementRepo := database.NewJobRequirementRepository(db)

	ctx := context.Background()

	// Define job requirements to seed
	jobRequirements := []models.JobRequirement{
		{
			Name:        "White Card",
			Description: stringPtr("Construction industry safety card"),
			IsActive:    true,
		},
		{
			Name:        "First Aid",
			Description: stringPtr("First aid certification required"),
			IsActive:    true,
		},
		{
			Name:        "Driver License",
			Description: stringPtr("Valid driver's license required"),
			IsActive:    true,
		},
		{
			Name:        "Own Tools",
			Description: stringPtr("Must provide own tools"),
			IsActive:    true,
		},
		{
			Name:        "Safety Boots",
			Description: stringPtr("Safety boots required"),
			IsActive:    true,
		},
		{
			Name:        "Hard Hat",
			Description: stringPtr("Hard hat required"),
			IsActive:    true,
		},
		{
			Name:        "High Vis Vest",
			Description: stringPtr("High visibility vest required"),
			IsActive:    true,
		},
		{
			Name:        "Experience Required",
			Description: stringPtr("Previous experience in the field required"),
			IsActive:    true,
		},
	}

	// Seed job requirements
	for _, requirement := range jobRequirements {
		// Check if requirement already exists
		existing, err := jobRequirementRepo.GetAll(ctx)
		if err == nil {
			// Check if this specific requirement already exists
			found := false
			for _, existingReq := range existing {
				if existingReq.Name == requirement.Name {
					fmt.Printf("Job requirement '%s' already exists, skipping...\n", requirement.Name)
					found = true
					break
				}
			}
			if found {
				continue
			}
		}

		// Create requirement using direct DB insert since we don't have Create method
		if err := db.WithContext(ctx).Create(&requirement).Error; err != nil {
			log.Printf("Failed to create job requirement '%s': %v", requirement.Name, err)
			continue
		}

		fmt.Printf("Successfully created job requirement '%s'\n", requirement.Name)
	}

	fmt.Println("Job requirements seeding completed!")
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
