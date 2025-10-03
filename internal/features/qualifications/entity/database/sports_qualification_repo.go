package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/qualifications/models"
	"gorm.io/gorm"
)

type SportsQualificationRepository interface {
	Create(ctx context.Context, sport *models.SportsQualification) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SportsQualification, error)
	GetByName(ctx context.Context, name string) (*models.SportsQualification, error)
	GetAll(ctx context.Context) ([]*models.SportsQualification, error)
	Update(ctx context.Context, sport *models.SportsQualification) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type sportsQualificationRepository struct {
	db *gorm.DB
}

func NewSportsQualificationRepository(db *gorm.DB) SportsQualificationRepository {
	return &sportsQualificationRepository{db: db}
}

func (r *sportsQualificationRepository) Create(ctx context.Context, sport *models.SportsQualification) error {
	return r.db.WithContext(ctx).Create(sport).Error
}

func (r *sportsQualificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SportsQualification, error) {
	var sport models.SportsQualification
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sport).Error
	if err != nil {
		return nil, err
	}
	return &sport, nil
}

func (r *sportsQualificationRepository) GetByName(ctx context.Context, name string) (*models.SportsQualification, error) {
	var sport models.SportsQualification
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&sport).Error
	if err != nil {
		return nil, err
	}
	return &sport, nil
}

func (r *sportsQualificationRepository) GetAll(ctx context.Context) ([]*models.SportsQualification, error) {
	var sports []*models.SportsQualification
	err := r.db.WithContext(ctx).Find(&sports).Error
	if err != nil {
		return nil, err
	}
	return sports, nil
}

func (r *sportsQualificationRepository) Update(ctx context.Context, sport *models.SportsQualification) error {
	return r.db.WithContext(ctx).Save(sport).Error
}

func (r *sportsQualificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.SportsQualification{}).Error
}
