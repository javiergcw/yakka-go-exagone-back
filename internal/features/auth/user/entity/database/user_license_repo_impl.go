package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/auth/user/models"
	"gorm.io/gorm"
)

type userLicenseRepository struct {
	db *gorm.DB
}

// NewUserLicenseRepository creates a new user license repository
func NewUserLicenseRepository(db *gorm.DB) UserLicenseRepository {
	return &userLicenseRepository{db: db}
}

func (r *userLicenseRepository) Create(ctx context.Context, license *models.UserLicense) error {
	return r.db.WithContext(ctx).Create(license).Error
}

func (r *userLicenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserLicense, error) {
	var license models.UserLicense
	err := r.db.WithContext(ctx).Preload("User").Preload("License").First(&license, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *userLicenseRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserLicense, error) {
	var licenses []*models.UserLicense
	err := r.db.WithContext(ctx).Preload("User").Preload("License").Where("user_id = ?", userID).Find(&licenses).Error
	if err != nil {
		return nil, err
	}
	return licenses, nil
}

func (r *userLicenseRepository) Update(ctx context.Context, license *models.UserLicense) error {
	return r.db.WithContext(ctx).Save(license).Error
}

func (r *userLicenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserLicense{}, "id = ?", id).Error
}

func (r *userLicenseRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserLicense{}, "user_id = ?", userID).Error
}

func (r *userLicenseRepository) CreateBatch(ctx context.Context, licenses []*models.UserLicense) error {
	return r.db.WithContext(ctx).CreateInBatches(licenses, 100).Error
}
