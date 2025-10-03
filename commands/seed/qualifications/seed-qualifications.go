package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
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

	// Create repositories
	sportsRepo := database.NewSportsQualificationRepository(db)
	qualificationRepo := database.NewQualificationRepository(db)

	ctx := context.Background()

	// Get all sports first
	sports, err := sportsRepo.GetAll(ctx)
	if err != nil {
		log.Fatalf("Failed to get sports: %v", err)
	}

	// Create sport-to-ID mapping
	sportMap := make(map[string]string)
	for _, sport := range sports {
		sportMap[sport.Name] = sport.ID.String()
	}

	// Define qualifications data
	qualificationsData := []struct {
		Sport          string
		Qualifications []struct {
			Title        string
			Organization string
			Country      string
		}
	}{
		{
			Sport: "Football (Soccer)",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"MiniRoos Certificate", "Football Australia", "Australia"},
				{"Grassroots Certificate", "Football Australia", "Australia"},
				{"Community Junior Licence", "Football Australia", "Australia"},
				{"Foundation of Football", "Football Australia", "Australia"},
				{"Skills Training Certificate", "Football Australia", "Australia"},
				{"Game Training Certificate", "Football Australia", "Australia"},
				{"Senior Coaching Certificate", "Football Australia", "Australia"},
				{"Community Courses", "Football Australia", "Australia"},
				{"UEFA A Diploma", "UEFA", "International"},
				{"UEFA B Diploma", "UEFA", "International"},
				{"UEFA C Diploma", "UEFA", "International"},
				{"UEFA PRO Diploma", "UEFA", "International"},
				{"UEFA Goalkeeping Level 1", "UEFA", "International"},
				{"UEFA Goalkeeping Level 2", "UEFA", "International"},
				{"UEFA Goalkeeping Level 3", "UEFA", "International"},
				{"AFC C Diploma", "AFC", "Asia"},
				{"AFC B Diploma", "AFC", "Asia"},
				{"AFC A Diploma", "AFC", "Asia"},
				{"AFC Pro Diploma", "AFC", "Asia"},
				{"AFC Goalkeeping Level 1", "AFC", "Asia"},
				{"AFC Goalkeeping Level 2", "AFC", "Asia"},
				{"AFC Goalkeeping Level 3", "AFC", "Asia"},
				{"Futsal Licence (Football Australia pathway)", "Football Australia", "Australia"},
			},
		},
		{
			Sport: "Rugby Union",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Smart Rugby (mandatory)", "Rugby Australia", "Australia"},
				{"Foundation Coach", "Rugby Australia", "Australia"},
				{"Development Coach", "Rugby Australia", "Australia"},
				{"Advanced Coach", "Rugby Australia", "Australia"},
				{"Performance Coach", "Rugby Australia", "Australia"},
			},
		},
		{
			Sport: "Rugby League",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Community Coach Accreditation", "Rugby League Australia", "Australia"},
				{"Club Coach", "Rugby League Australia", "Australia"},
				{"Performance Coach", "Rugby League Australia", "Australia"},
			},
		},
		{
			Sport: "Australian Rules Football (AFL)",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "AFL", "Australia"},
				{"Development Coach", "AFL", "Australia"},
				{"Performance Coach", "AFL", "Australia"},
			},
		},
		{
			Sport: "Touch Football",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Referee Accreditation Level 1", "Touch Football Australia", "Australia"},
				{"Referee Accreditation Level 2", "Touch Football Australia", "Australia"},
				{"Referee Accreditation Level 3", "Touch Football Australia", "Australia"},
			},
		},
		{
			Sport: "American Football (Gridiron)",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Gridiron Australia Coach Accreditation", "Gridiron Australia", "Australia"},
			},
		},
		{
			Sport: "Flag Football",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Flag Football Coach Accreditation", "Flag Football Association", "International"},
			},
		},
		{
			Sport: "Tennis",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Participation (Foundation) Coaching Course", "Tennis Australia", "Australia"},
				{"Development (Level 1) Coaching", "Tennis Australia", "Australia"},
				{"Performance/High Performance (Level 3)", "Tennis Australia", "Australia"},
			},
		},
		{
			Sport: "Badminton",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Badminton Australia", "Australia"},
				{"Development Coach", "Badminton Australia", "Australia"},
				{"Advanced Coach", "Badminton Australia", "Australia"},
			},
		},
		{
			Sport: "Squash",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Squash Australia", "Australia"},
				{"Development Coach", "Squash Australia", "Australia"},
				{"Advanced Coach", "Squash Australia", "Australia"},
			},
		},
		{
			Sport: "Table Tennis",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Table Tennis Australia", "Australia"},
				{"Development Coach", "Table Tennis Australia", "Australia"},
				{"Advanced Coach", "Table Tennis Australia", "Australia"},
			},
		},
		{
			Sport: "Pickleball",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Deliverer Accreditation", "Pickleball Association", "International"},
				{"Coach Accreditation", "Pickleball Association", "International"},
			},
		},
		{
			Sport: "Padel",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Deliverer Accreditation", "Padel Federation", "International"},
				{"Coach Accreditation", "Padel Federation", "International"},
			},
		},
		{
			Sport: "Swimming",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"AUSTSWIM Teacher of Swimming & Water Safety (TSW)", "AUSTSWIM", "Australia"},
				{"Development Coach", "Swimming Australia", "Australia"},
				{"Advanced Coach", "Swimming Australia", "Australia"},
				{"Performance Coach", "Swimming Australia", "Australia"},
				{"Pool Lifeguard", "Swimming Australia", "Australia"},
			},
		},
		{
			Sport: "Surf Life Saving",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Bronze Medallion", "Surf Life Saving Australia", "Australia"},
			},
		},
		{
			Sport: "Water Polo",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Water Polo Australia", "Australia"},
				{"Development Coach", "Water Polo Australia", "Australia"},
				{"Performance Coach", "Water Polo Australia", "Australia"},
			},
		},
		{
			Sport: "Surfing",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Surfing Australia", "Australia"},
				{"Level 1", "Surfing Australia", "Australia"},
				{"Level 2", "Surfing Australia", "Australia"},
				{"High Performance", "Surfing Australia", "Australia"},
			},
		},
		{
			Sport: "Sailing",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Dinghy Instructor", "Australian Sailing", "Australia"},
				{"Powerboat Handling", "Australian Sailing", "Australia"},
				{"Safety Boat Operator", "Australian Sailing", "Australia"},
			},
		},
		{
			Sport: "Windsurfing",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Assistant Instructor", "Australian Windsurfing", "Australia"},
				{"Instructor", "Australian Windsurfing", "Australia"},
			},
		},
		{
			Sport: "Rowing",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Coach", "Rowing Australia", "Australia"},
				{"Level 2 Coach", "Rowing Australia", "Australia"},
				{"Level 3 Coach", "Rowing Australia", "Australia"},
				{"Level 4 Coach", "Rowing Australia", "Australia"},
			},
		},
		{
			Sport: "Canoe/Kayak/Surf Ski",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Instructor Qualifications", "Australian Canoeing", "Australia"},
				{"Foundation Coach", "Australian Canoeing", "Australia"},
			},
		},
		{
			Sport: "Kite Surfing",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Instructor (ITC)", "International Kiteboarding Association", "International"},
			},
		},
		{
			Sport: "Boxing",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Bronze Coach", "Boxing Australia", "Australia"},
				{"Silver Coach", "Boxing Australia", "Australia"},
				{"Gold Coach", "Boxing Australia", "Australia"},
			},
		},
		{
			Sport: "Judo",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Assistant Coach", "Judo Australia", "Australia"},
				{"Sensei Coach", "Judo Australia", "Australia"},
				{"Senior Coach", "Judo Australia", "Australia"},
				{"Advanced Senior Coach", "Judo Australia", "Australia"},
			},
		},
		{
			Sport: "Karate",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Accredited Karate Instructor", "Karate Federation", "International"},
			},
		},
		{
			Sport: "Taekwondo",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Accredited Coach (Kyorugi/Poomsae)", "Taekwondo Federation", "International"},
			},
		},
		{
			Sport: "Wrestling",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Coach (with WAL requirements)", "Wrestling Australia", "Australia"},
			},
		},
		{
			Sport: "Brazilian Jiu-Jitsu",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"White Belt", "IBJJF", "International"},
				{"Blue Belt", "IBJJF", "International"},
				{"Purple Belt", "IBJJF", "International"},
				{"Brown Belt", "IBJJF", "International"},
				{"Black Belt (+ degrees)", "IBJJF", "International"},
				{"Rules/Referee Course", "IBJJF", "International"},
			},
		},
		{
			Sport: "Athletics",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Youth", "Athletics Australia", "Australia"},
				{"Development Coach", "Athletics Australia", "Australia"},
				{"Performance Coach", "Athletics Australia", "Australia"},
				{"Level 4 High Performance", "Athletics Australia", "Australia"},
			},
		},
		{
			Sport: "Triathlon/Duathlon",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Triathlon Australia", "Australia"},
				{"Development Coach", "Triathlon Australia", "Australia"},
				{"Performance Coach", "Triathlon Australia", "Australia"},
				{"High Performance Coach", "Triathlon Australia", "Australia"},
			},
		},
		{
			Sport: "Gymnastics",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Fundamental Coach", "Gymnastics Australia", "Australia"},
				{"Beginner Coach", "Gymnastics Australia", "Australia"},
				{"Intermediate Coach", "Gymnastics Australia", "Australia"},
				{"Advanced Coach", "Gymnastics Australia", "Australia"},
				{"High Performance Coach", "Gymnastics Australia", "Australia"},
			},
		},
		{
			Sport: "Weightlifting",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Club Coach", "Weightlifting Australia", "Australia"},
				{"Level 2 State Coach", "Weightlifting Australia", "Australia"},
				{"Level 3 National Coach", "Weightlifting Australia", "Australia"},
			},
		},
		{
			Sport: "Powerlifting",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"PA NCAS Level 1 Coach (and above)", "Powerlifting Australia", "Australia"},
			},
		},
		{
			Sport: "Strength & Conditioning",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"ASCA Level 1 Coach", "ASCA", "Australia"},
				{"ASCA Level 2 Coach", "ASCA", "Australia"},
				{"ASCA Level 3 Coach", "ASCA", "Australia"},
			},
		},
		{
			Sport: "Netball",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Netball Australia", "Australia"},
				{"Development Coach", "Netball Australia", "Australia"},
				{"Intermediate Coach", "Netball Australia", "Australia"},
				{"Advanced Coach", "Netball Australia", "Australia"},
				{"Elite Coach", "Netball Australia", "Australia"},
				{"High Performance Coach", "Netball Australia", "Australia"},
				{"Umpire C Badge", "Netball Australia", "Australia"},
				{"Umpire B Badge", "Netball Australia", "Australia"},
				{"Umpire A Badge", "Netball Australia", "Australia"},
				{"Umpire AA Badge", "Netball Australia", "Australia"},
			},
		},
		{
			Sport: "Volleyball",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Coach", "Volleyball Australia", "Australia"},
				{"Level 2 Coach", "Volleyball Australia", "Australia"},
				{"Level 3 Coach", "Volleyball Australia", "Australia"},
				{"Level 4 Coach", "Volleyball Australia", "Australia"},
			},
		},
		{
			Sport: "Softball",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Softball Australia", "Australia"},
				{"Performance Talent Coach", "Softball Australia", "Australia"},
				{"High Performance/Elite Coach", "Softball Australia", "Australia"},
				{"Master Coach", "Softball Australia", "Australia"},
			},
		},
		{
			Sport: "Baseball",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Certification A", "Baseball Australia", "Australia"},
				{"Certification B", "Baseball Australia", "Australia"},
				{"Certification C", "Baseball Australia", "Australia"},
			},
		},
		{
			Sport: "Cricket",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Community Coach", "Cricket Australia", "Australia"},
				{"Advanced Coach", "Cricket Australia", "Australia"},
				{"Representative Coach", "Cricket Australia", "Australia"},
				{"High Performance Coach", "Cricket Australia", "Australia"},
			},
		},
		{
			Sport: "Hockey (Field)",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Foundation Coach", "Hockey Australia", "Australia"},
				{"Development Coach", "Hockey Australia", "Australia"},
				{"Advanced Coach (+RCC)", "Hockey Australia", "Australia"},
			},
		},
		{
			Sport: "Ice Hockey",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"NCACS Level 1", "Ice Hockey Australia", "Australia"},
				{"NCACS Level 2", "Ice Hockey Australia", "Australia"},
				{"NCACS Level 3", "Ice Hockey Australia", "Australia"},
				{"NCACS Level 4", "Ice Hockey Australia", "Australia"},
			},
		},
		{
			Sport: "Skiing/Snowboarding",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"APSI Instructor – Alpine", "APSI", "Australia"},
				{"APSI Instructor – Snowboard", "APSI", "Australia"},
				{"APSI Instructor – Telemark", "APSI", "Australia"},
				{"APSI Instructor – Nordic", "APSI", "Australia"},
			},
		},
		{
			Sport: "Ice Skating",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"APSA Level 0", "APSA", "Australia"},
				{"APSA Level 1", "APSA", "Australia"},
				{"APSA Level 2", "APSA", "Australia"},
				{"APSA Level 3+", "APSA", "Australia"},
			},
		},
		{
			Sport: "Curling",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Coaching Instructor", "Curling Association", "International"},
				{"Competition Coach", "Curling Association", "International"},
				{"Performance Coach", "Curling Association", "International"},
			},
		},
		{
			Sport: "Cycling",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Community Coach", "Cycling Australia", "Australia"},
				{"Foundation Coach", "Cycling Australia", "Australia"},
				{"Development Coach", "Cycling Australia", "Australia"},
				{"Advanced Coach", "Cycling Australia", "Australia"},
				{"Elite Coach", "Cycling Australia", "Australia"},
			},
		},
		{
			Sport: "Equestrian",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Introductory Coach", "Equestrian Australia", "Australia"},
				{"Level 1 Coach", "Equestrian Australia", "Australia"},
				{"Level 2 Coach", "Equestrian Australia", "Australia"},
				{"Level 3 Coach", "Equestrian Australia", "Australia"},
			},
		},
		{
			Sport: "Archery",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Community Coach", "Archery Australia", "Australia"},
				{"Level 1 Coach", "Archery Australia", "Australia"},
				{"Level 2 Coach", "Archery Australia", "Australia"},
				{"Level 3 Coach", "Archery Australia", "Australia"},
				{"Level 4 Coach", "Archery Australia", "Australia"},
			},
		},
		{
			Sport: "Shooting",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Club Coach", "Shooting Australia", "Australia"},
				{"Competition Coach", "Shooting Australia", "Australia"},
				{"Advanced Coach", "Shooting Australia", "Australia"},
			},
		},
		{
			Sport: "Lawn Bowls",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Introductory Coach", "Bowls Australia", "Australia"},
				{"Club Coach", "Bowls Australia", "Australia"},
				{"Advanced Coach", "Bowls Australia", "Australia"},
				{"High Performance Coach", "Bowls Australia", "Australia"},
			},
		},
		{
			Sport: "Golf",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Community Program Deliverer", "Golf Australia", "Australia"},
				{"National Program Deliverer (MyGolf/Get Into Golf)", "Golf Australia", "Australia"},
			},
		},
		{
			Sport: "CrossFit",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Certificate", "CrossFit", "International"},
				{"Level 2 Certificate", "CrossFit", "International"},
				{"Level 3 CCFT", "CrossFit", "International"},
				{"Level 4 CF-L4", "CrossFit", "International"},
			},
		},
		{
			Sport: "Yoga",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Level 1 Teacher (350 hrs)", "Yoga Alliance", "International"},
				{"Level 2 Teacher (500 hrs)", "Yoga Alliance", "International"},
				{"Level 3 Teacher (1000 hrs)", "Yoga Alliance", "International"},
			},
		},
		{
			Sport: "Pilates",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"Certificate IV Instructor", "Pilates Alliance", "International"},
				{"Diploma Instructor (PAA recognised)", "Pilates Alliance", "International"},
			},
		},
		{
			Sport: "General Sports",
			Qualifications: []struct {
				Title        string
				Organization string
				Country      string
			}{
				{"SIS20122 – Certificate II in Sport and Recreation", "Australian SIS", "Australia"},
				{"SIS20321 – Certificate II in Sport Coaching", "Australian SIS", "Australia"},
				{"SIS30122 – Certificate III in Sport, Aquatics and Recreation", "Australian SIS", "Australia"},
				{"SIS30521 – Certificate III in Sport Coaching", "Australian SIS", "Australia"},
				{"SIS40321 – Certificate IV in Sport Coaching", "Australian SIS", "Australia"},
				{"SIS50321 – Diploma of Sport", "Australian SIS", "Australia"},
				{"SIS50122 – Diploma of Sport, Aquatics and Recreation Management", "Australian SIS", "Australia"},
				{"Diploma of Sport & Exercise Science (provider specific)", "University Programs", "International"},
				{"Doctoral Degree (PhD) in Sport", "University Programs", "International"},
				{"Doctoral Degree (PhD) Exercise", "University Programs", "International"},
				{"Doctoral Degree (PhD) Science", "University Programs", "International"},
				{"Master's Degree in Sport", "University Programs", "International"},
				{"Master's Degree in Exercise Science", "University Programs", "International"},
				{"Master's Degree in Sport Management", "University Programs", "International"},
			},
		},
	}

	// Seed qualifications
	for _, sportData := range qualificationsData {
		sportID, exists := sportMap[sportData.Sport]
		if !exists {
			fmt.Printf("Sport '%s' not found, skipping qualifications...\n", sportData.Sport)
			continue
		}

		// Parse sport ID
		sportUUID, err := uuid.Parse(sportID)
		if err != nil {
			fmt.Printf("Invalid sport ID for '%s': %v\n", sportData.Sport, err)
			continue
		}

		fmt.Printf("Processing qualifications for sport: %s\n", sportData.Sport)

		for _, q := range sportData.Qualifications {
			qualification := models.Qualification{
				SportID:      sportUUID,
				Title:        q.Title,
				Organization: stringPtr(q.Organization),
				Country:      stringPtr(q.Country),
				Status:       "active",
			}

			// Create qualification
			if err := qualificationRepo.Create(ctx, &qualification); err != nil {
				fmt.Printf("Failed to create qualification '%s' for sport '%s': %v\n", q.Title, sportData.Sport, err)
				continue
			}

			fmt.Printf("Successfully created qualification: %s (%s)\n", q.Title, sportData.Sport)
		}
	}

	fmt.Println("Qualifications seeding completed!")
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
