package rest

import (
	"context"
	"net/http"

	"github.com/yakka-backend/internal/features/qualifications/entity/database"
	"github.com/yakka-backend/internal/features/qualifications/models"
	"github.com/yakka-backend/internal/shared/response"
)

// QualificationHandler handles qualification-related HTTP requests
type QualificationHandler struct {
	qualificationRepo database.QualificationRepository
}

// NewQualificationHandler creates a new qualification handler
func NewQualificationHandler(qualificationRepo database.QualificationRepository) *QualificationHandler {
	return &QualificationHandler{
		qualificationRepo: qualificationRepo,
	}
}

// QualificationItem represents a single qualification
type QualificationItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Organization string `json:"organization,omitempty"`
	Country      string `json:"country,omitempty"`
	Status       string `json:"status"`
}

// SportWithQualifications represents a sport with its qualifications
type SportWithQualifications struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Qualifications []QualificationItem `json:"qualifications"`
}

// GetQualificationsResponse represents the response for getting all qualifications
type GetQualificationsResponse struct {
	Sports  []SportWithQualifications `json:"sports"`
	Total   int                       `json:"total"`
	Message string                    `json:"message"`
}

// GetQualifications handles GET /api/v1/qualifications
func (h *QualificationHandler) GetQualifications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get all qualifications with sport information
	qualifications, err := h.qualificationRepo.GetAllWithSport(ctx)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get qualifications")
		return
	}

	// Group qualifications by sport
	sportMap := make(map[string]*SportWithQualifications)

	for _, q := range qualifications {
		if q.Sport == nil {
			continue
		}

		sportID := q.Sport.ID.String()
		sportName := q.Sport.Name

		// Create sport entry if it doesn't exist
		if sportMap[sportID] == nil {
			sportMap[sportID] = &SportWithQualifications{
				ID:             sportID,
				Name:           sportName,
				Qualifications: []QualificationItem{},
			}
		}

		// Add qualification to the sport
		qualification := QualificationItem{
			ID:           q.ID.String(),
			Title:        q.Title,
			Organization: getStringValue(q.Organization),
			Country:      getStringValue(q.Country),
			Status:       q.Status,
		}

		sportMap[sportID].Qualifications = append(sportMap[sportID].Qualifications, qualification)
	}

	// Convert map to slice
	sports := make([]SportWithQualifications, 0, len(sportMap))
	totalQualifications := 0
	for _, sport := range sportMap {
		sports = append(sports, *sport)
		totalQualifications += len(sport.Qualifications)
	}

	resp := GetQualificationsResponse{
		Sports:  sports,
		Total:   totalQualifications,
		Message: "Qualifications retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// Helper function to get sport name safely
func getSportName(sport *models.SportsQualification) string {
	if sport == nil {
		return ""
	}
	return sport.Name
}

// Helper function to get string value safely
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
