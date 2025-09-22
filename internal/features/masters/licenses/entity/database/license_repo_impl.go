package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/licenses/models"
	"gorm.io/gorm"
)

type licenseRepository struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) LicenseRepository {
	return &licenseRepository{db: db}
}

func (r *licenseRepository) Create(ctx context.Context, license *models.License) error {
	return r.db.WithContext(ctx).Create(license).Error
}

func (r *licenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.License, error) {
	var license models.License
	err := r.db.WithContext(ctx).First(&license, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *licenseRepository) GetAll(ctx context.Context) ([]*models.License, error) {
	var licenses []*models.License
	err := r.db.WithContext(ctx).Find(&licenses).Error
	if err != nil {
		return nil, err
	}
	return licenses, nil
}

func (r *licenseRepository) Update(ctx context.Context, license *models.License) error {
	return r.db.WithContext(ctx).Save(license).Error
}

func (r *licenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.License{}, "id = ?", id).Error
}
