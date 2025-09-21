package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/builder_profiles/models"
	"gorm.io/gorm"
)

// builderProfileRepository implements BuilderProfileRepository
type builderProfileRepository struct {
	db *gorm.DB
}

// NewBuilderProfileRepository creates a new builder profile repository
func NewBuilderProfileRepository(db *gorm.DB) BuilderProfileRepository {
	return &builderProfileRepository{
		db: db,
	}
}

// Create creates a new builder profile
func (r *builderProfileRepository) Create(ctx context.Context, profile *models.BuilderProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByID retrieves a builder profile by ID
func (r *builderProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BuilderProfile, error) {
	var profile models.BuilderProfile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetByUserID retrieves a builder profile by user ID
func (r *builderProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.BuilderProfile, error) {
	var profile models.BuilderProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update updates a builder profile
func (r *builderProfileRepository) Update(ctx context.Context, profile *models.BuilderProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Delete deletes a builder profile by ID
func (r *builderProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BuilderProfile{}).Error
}

// DeleteByUserID deletes a builder profile by user ID
func (r *builderProfileRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.BuilderProfile{}).Error
}
