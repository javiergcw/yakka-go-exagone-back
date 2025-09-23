package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SkillSubcategory represents a skill subcategory in the system
type SkillSubcategory struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryID  uuid.UUID      `json:"category_id" gorm:"type:uuid;not null"`
	Name        string         `json:"name" gorm:"not null;size:255"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null;type:timestamptz"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Category SkillCategory `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName returns the table name for the SkillSubcategory model
func (SkillSubcategory) TableName() string {
	return "skill_subcategories"
}
