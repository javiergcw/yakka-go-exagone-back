package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yakka-backend/internal/features/qualifications/entity/database"
	"github.com/yakka-backend/internal/features/qualifications/models"
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
	sportsRepo := database.NewSportsQualificationRepository(db)

	ctx := context.Background()

	// Define sports to seed
	sports := []models.SportsQualification{
		{Name: "Football (Soccer)"},
		{Name: "Rugby Union"},
		{Name: "Rugby League"},
		{Name: "Australian Rules Football (AFL)"},
		{Name: "Touch Football"},
		{Name: "American Football (Gridiron)"},
		{Name: "Flag Football"},
		{Name: "Tennis"},
		{Name: "Badminton"},
		{Name: "Squash"},
		{Name: "Table Tennis"},
		{Name: "Pickleball"},
		{Name: "Padel"},
		//
		{Name: "Swimming"},
		{Name: "Surf Life Saving"},
		{Name: "Water Polo"},
		{Name: "Surfing"},
		{Name: "Sailing"},
		{Name: "Windsurfing"},
		{Name: "Rowing"},
		{Name: "Canoe/Kayak/Surf Ski"},
		{Name: "Kite Surfing"},
		{Name: "Boxing"},
		{Name: "Judo"},
		{Name: "Karate"},
		{Name: "Taekwondo"},
		{Name: "Wrestling"},
		{Name: "Brazilian Jiu-Jitsu"},
		{Name: "Athletics"},
		{Name: "Triathlon/Duathlon"},
		{Name: "Gymnastics"},
		{Name: "Weightlifting"},
		{Name: "Powerlifting"},
		{Name: "Strength & Conditioning"},
		{Name: "Netball"},
		{Name: "Volleyball"},
		{Name: "Softball"},
		{Name: "Baseball"},
		{Name: "Cricket"},
		{Name: "Hockey (Field)"},
		{Name: "Ice Hockey"},
		{Name: "Skiing/Snowboarding"},
		{Name: "Ice Skating"},
		{Name: "Curling"},
		{Name: "Cycling"},
		{Name: "Equestrian"},
		{Name: "Archery"},
		{Name: "Shooting"},
		{Name: "Lawn Bowls"},
		{Name: "Golf"},
		{Name: "CrossFit"},
		{Name: "Yoga"},
		{Name: "Pilates"},
		{Name: "General Sports"},
	}

	// Seed sports
	for _, sport := range sports {
		// Check if sport already exists
		existing, err := sportsRepo.GetByName(ctx, sport.Name)
		if err == nil && existing != nil {
			fmt.Printf("Sport '%s' already exists, skipping...\n", sport.Name)
			continue
		}

		// Create sport
		if err := sportsRepo.Create(ctx, &sport); err != nil {
			log.Printf("Failed to create sport '%s': %v", sport.Name, err)
			continue
		}

		fmt.Printf("Successfully created sport '%s'\n", sport.Name)
	}

	fmt.Println("Sports qualifications seeding completed!")
}
