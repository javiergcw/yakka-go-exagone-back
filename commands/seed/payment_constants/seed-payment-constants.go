package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yakka-backend/internal/features/masters/payment_constants/entity/database"
	"github.com/yakka-backend/internal/features/masters/payment_constants/models"
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
	paymentConstantRepo := database.NewPaymentConstantRepository(db)

	ctx := context.Background()

	// Define payment constants to seed
	paymentConstants := []models.PaymentConstant{
		{
			Name:        "GST",
			Value:       10, // 10% GST
			Description: stringPtr("Goods and Services Tax percentage"),
			IsActive:    true,
		},
		{
			Name:        "WAGE_HOURLY",
			Value:       28, // $28 per hour
			Description: stringPtr("Default hourly wage rate"),
			IsActive:    true,
		},
	}

	// Seed payment constants
	for _, constant := range paymentConstants {
		// Check if constant already exists
		existing, err := paymentConstantRepo.GetByName(ctx, constant.Name)
		if err == nil && existing != nil {
			fmt.Printf("Payment constant '%s' already exists, skipping...\n", constant.Name)
			continue
		}

		// Create constant
		if err := paymentConstantRepo.Create(ctx, &constant); err != nil {
			log.Printf("Failed to create payment constant '%s': %v", constant.Name, err)
			continue
		}

		fmt.Printf("Successfully created payment constant '%s' with value %d\n", constant.Name, constant.Value)
	}

	fmt.Println("Payment constants seeding completed!")
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
