package models

import (
	"time"

	"github.com/google/uuid"
	experience_models "github.com/yakka-backend/internal/features/masters/experience_levels/models"
	skill_models "github.com/yakka-backend/internal/features/masters/skills/models"
)

// LabourProfileSkill represents a skill associated with a labour profile
type LabourProfileSkill struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LabourProfileID   uuid.UUID `json:"labour_profile_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_labour_profile_subcategory"`
	CategoryID        uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
	SubcategoryID     uuid.UUID `json:"subcategory_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_labour_profile_subcategory"`
	ExperienceLevelID uuid.UUID `json:"experience_level_id" gorm:"type:uuid;not null;index"`
	YearsExperience   float64   `json:"years_experience" gorm:"type:decimal(4,1)"`
	IsPrimary         bool      `json:"is_primary" gorm:"not null;default:false"`
	CreatedAt         time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`

	// Relationships
	LabourProfile   LabourProfile                     `json:"labour_profile" gorm:"foreignKey:LabourProfileID;references:ID"`
	Category        skill_models.SkillCategory        `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	Subcategory     skill_models.SkillSubcategory     `json:"subcategory" gorm:"foreignKey:SubcategoryID;references:ID"`
	ExperienceLevel experience_models.ExperienceLevel `json:"experience_level" gorm:"foreignKey:ExperienceLevelID;references:ID"`
}

// TableName returns the table name for the LabourProfileSkill model
func (LabourProfileSkill) TableName() string {
	return "labour_profile_skills"
}

// BeforeCreate sets up unique constraint
func (lps *LabourProfileSkill) BeforeCreate() error {
	// GORM will handle the unique constraint through the index definition
	return nil
}
