package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/jobs/models"
)

// JobLicenseRepository defines the interface for job license operations
type JobLicenseRepository interface {
	Create(ctx context.Context, jobLicense *models.JobLicense) error
	GetByJobID(ctx context.Context, jobID uuid.UUID) ([]*models.JobLicense, error)
	DeleteByJobID(ctx context.Context, jobID uuid.UUID) error
	DeleteByJobAndLicense(ctx context.Context, jobID, licenseID uuid.UUID) error
}
