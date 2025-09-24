package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
)

// JobSkillRepository defines the interface for job skill operations
type JobSkillRepository interface {
	Create(ctx context.Context, jobSkill *models.JobSkill) error
	GetByJobID(ctx context.Context, jobID uuid.UUID) ([]*models.JobSkill, error)
	DeleteByJobID(ctx context.Context, jobID uuid.UUID) error
	DeleteByJobAndSkill(ctx context.Context, jobID, skillCategoryID, skillSubcategoryID uuid.UUID) error
}
