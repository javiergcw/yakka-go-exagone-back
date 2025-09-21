package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/labour_profiles/models"
	"gorm.io/gorm"
)

// labourProfileRepository implements LabourProfileRepository
type labourProfileRepository struct {
	db *gorm.DB
}

// NewLabourProfileRepository creates a new labour profile repository
func NewLabourProfileRepository(db *gorm.DB) LabourProfileRepository {
	return &labourProfileRepository{
		db: db,
	}
}

// Create creates a new labour profile
func (r *labourProfileRepository) Create(ctx context.Context, profile *models.LabourProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByID retrieves a labour profile by ID
func (r *labourProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LabourProfile, error) {
	var profile models.LabourProfile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetByUserID retrieves a labour profile by user ID
func (r *labourProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.LabourProfile, error) {
	var profile models.LabourProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update updates a labour profile
func (r *labourProfileRepository) Update(ctx context.Context, profile *models.LabourProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Delete deletes a labour profile by ID
func (r *labourProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.LabourProfile{}).Error
}

// DeleteByUserID deletes a labour profile by user ID
func (r *labourProfileRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.LabourProfile{}).Error
}
