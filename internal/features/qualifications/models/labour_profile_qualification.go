package models

import (
	"time"

	"github.com/google/uuid"
)

// LabourProfileQualification represents the relationship between a labour profile and a qualification
type LabourProfileQualification struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LabourProfileID uuid.UUID      `json:"labour_profile_id" gorm:"type:uuid;not null;index"`
	QualificationID uuid.UUID      `json:"qualification_id" gorm:"type:uuid;not null;index"`
	Qualification   *Qualification `json:"qualification,omitempty" gorm:"foreignKey:QualificationID"`
	DateObtained    *time.Time     `json:"date_obtained" gorm:"type:date"`
	ExpiresAt       *time.Time     `json:"expires_at" gorm:"type:date"`
	Status          string         `json:"status" gorm:"size:50;default:valid"`
	CreatedAt       time.Time      `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the LabourProfileQualification model
func (LabourProfileQualification) TableName() string {
	return "labour_profile_qualifications"
}
