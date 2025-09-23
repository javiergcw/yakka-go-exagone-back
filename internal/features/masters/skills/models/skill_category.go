package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SkillCategory represents a skill category in the system
type SkillCategory struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null;size:255;unique"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null;type:timestamptz"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for the SkillCategory model
func (SkillCategory) TableName() string {
	return "skill_categories"
}
