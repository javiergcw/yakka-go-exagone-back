package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	"github.com/yakka-backend/internal/features/masters/job_types/models"
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
	jobTypeRepo := database.NewJobTypeRepository(db)

	ctx := context.Background()

	// Define job types to seed
	jobTypes := []models.JobType{
		{
			Name:        "Casual Job",
			Description: stringPtr("Casual employment with flexible hours"),
			IsActive:    true,
		},
		{
			Name:        "Part Time",
			Description: stringPtr("Part-time employment"),
			IsActive:    true,
		},
		{
			Name:        "Full Time",
			Description: stringPtr("Full-time employment"),
			IsActive:    true,
		},
		{
			Name:        "Farms Job",
			Description: stringPtr("Agricultural and farming work"),
			IsActive:    true,
		},
		{
			Name:        "Mining Job",
			Description: stringPtr("Mining industry work"),
			IsActive:    true,
		},
		{
			Name:        "FIFO",
			Description: stringPtr("Fly In Fly Out work arrangements"),
			IsActive:    true,
		},
		{
			Name:        "Seasonal Job",
			Description: stringPtr("Seasonal employment"),
			IsActive:    true,
		},
		{
			Name:        "W&H Visa",
			Description: stringPtr("Working Holiday Visa jobs"),
			IsActive:    true,
		},
		{
			Name:        "Other",
			Description: stringPtr("Other types of employment"),
			IsActive:    true,
		},
	}

	// Seed job types
	for _, jobType := range jobTypes {
		// Check if job type already exists
		existing, err := jobTypeRepo.GetAll(ctx)
		if err == nil {
			// Check if this specific job type already exists
			found := false
			for _, existingType := range existing {
				if existingType.Name == jobType.Name {
					fmt.Printf("Job type '%s' already exists, skipping...\n", jobType.Name)
					found = true
					break
				}
			}
			if found {
				continue
			}
		}

		// Create job type using direct DB insert since we don't have Create method
		if err := db.WithContext(ctx).Create(&jobType).Error; err != nil {
			log.Printf("Failed to create job type '%s': %v", jobType.Name, err)
			continue
		}

		fmt.Printf("Successfully created job type '%s'\n", jobType.Name)
	}

	fmt.Println("Job types seeding completed!")
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
