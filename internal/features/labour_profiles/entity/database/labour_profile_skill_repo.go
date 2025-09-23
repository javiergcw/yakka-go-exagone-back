package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/labour_profiles/models"
)

// LabourProfileSkillRepository defines the interface for labour profile skill database operations
type LabourProfileSkillRepository interface {
	Create(ctx context.Context, skill *models.LabourProfileSkill) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LabourProfileSkill, error)
	GetByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) ([]*models.LabourProfileSkill, error)
	Update(ctx context.Context, skill *models.LabourProfileSkill) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) error
	CreateBatch(ctx context.Context, skills []*models.LabourProfileSkill) error
}
