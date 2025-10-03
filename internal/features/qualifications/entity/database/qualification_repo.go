package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/qualifications/models"
	"gorm.io/gorm"
)

type QualificationRepository interface {
	Create(ctx context.Context, qualification *models.Qualification) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Qualification, error)
	GetBySportID(ctx context.Context, sportID uuid.UUID) ([]*models.Qualification, error)
	GetAll(ctx context.Context) ([]*models.Qualification, error)
	GetAllWithSport(ctx context.Context) ([]*models.Qualification, error)
	Update(ctx context.Context, qualification *models.Qualification) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type qualificationRepository struct {
	db *gorm.DB
}

func NewQualificationRepository(db *gorm.DB) QualificationRepository {
	return &qualificationRepository{db: db}
}

func (r *qualificationRepository) Create(ctx context.Context, qualification *models.Qualification) error {
	return r.db.WithContext(ctx).Create(qualification).Error
}

func (r *qualificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Qualification, error) {
	var qualification models.Qualification
	err := r.db.WithContext(ctx).Preload("Sport").Where("id = ?", id).First(&qualification).Error
	if err != nil {
		return nil, err
	}
	return &qualification, nil
}

func (r *qualificationRepository) GetBySportID(ctx context.Context, sportID uuid.UUID) ([]*models.Qualification, error) {
	var qualifications []*models.Qualification
	err := r.db.WithContext(ctx).Preload("Sport").Where("sport_id = ?", sportID).Find(&qualifications).Error
	if err != nil {
		return nil, err
	}
	return qualifications, nil
}

func (r *qualificationRepository) GetAll(ctx context.Context) ([]*models.Qualification, error) {
	var qualifications []*models.Qualification
	err := r.db.WithContext(ctx).Find(&qualifications).Error
	if err != nil {
		return nil, err
	}
	return qualifications, nil
}

func (r *qualificationRepository) GetAllWithSport(ctx context.Context) ([]*models.Qualification, error) {
	var qualifications []*models.Qualification
	err := r.db.WithContext(ctx).Preload("Sport").Find(&qualifications).Error
	if err != nil {
		return nil, err
	}
	return qualifications, nil
}

func (r *qualificationRepository) Update(ctx context.Context, qualification *models.Qualification) error {
	return r.db.WithContext(ctx).Save(qualification).Error
}

func (r *qualificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Qualification{}).Error
}
