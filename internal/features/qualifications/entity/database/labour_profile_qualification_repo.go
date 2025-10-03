package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/qualifications/models"
	"gorm.io/gorm"
)

type LabourProfileQualificationRepository interface {
	Create(ctx context.Context, profileQualification *models.LabourProfileQualification) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LabourProfileQualification, error)
	GetByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) ([]*models.LabourProfileQualification, error)
	GetByQualificationID(ctx context.Context, qualificationID uuid.UUID) ([]*models.LabourProfileQualification, error)
	GetByLabourProfileAndQualification(ctx context.Context, labourProfileID, qualificationID uuid.UUID) (*models.LabourProfileQualification, error)
	Update(ctx context.Context, profileQualification *models.LabourProfileQualification) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) error
	DeleteByLabourProfileAndQualification(ctx context.Context, labourProfileID, qualificationID uuid.UUID) error
}

type labourProfileQualificationRepository struct {
	db *gorm.DB
}

func NewLabourProfileQualificationRepository(db *gorm.DB) LabourProfileQualificationRepository {
	return &labourProfileQualificationRepository{db: db}
}

func (r *labourProfileQualificationRepository) Create(ctx context.Context, profileQualification *models.LabourProfileQualification) error {
	return r.db.WithContext(ctx).Create(profileQualification).Error
}

func (r *labourProfileQualificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LabourProfileQualification, error) {
	var profileQualification models.LabourProfileQualification
	err := r.db.WithContext(ctx).Preload("Qualification.Sport").Where("id = ?", id).First(&profileQualification).Error
	if err != nil {
		return nil, err
	}
	return &profileQualification, nil
}

func (r *labourProfileQualificationRepository) GetByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) ([]*models.LabourProfileQualification, error) {
	var profileQualifications []*models.LabourProfileQualification
	err := r.db.WithContext(ctx).Preload("Qualification.Sport").Where("labour_profile_id = ?", labourProfileID).Find(&profileQualifications).Error
	if err != nil {
		return nil, err
	}
	return profileQualifications, nil
}

func (r *labourProfileQualificationRepository) GetByQualificationID(ctx context.Context, qualificationID uuid.UUID) ([]*models.LabourProfileQualification, error) {
	var profileQualifications []*models.LabourProfileQualification
	err := r.db.WithContext(ctx).Preload("Qualification.Sport").Where("qualification_id = ?", qualificationID).Find(&profileQualifications).Error
	if err != nil {
		return nil, err
	}
	return profileQualifications, nil
}

func (r *labourProfileQualificationRepository) GetByLabourProfileAndQualification(ctx context.Context, labourProfileID, qualificationID uuid.UUID) (*models.LabourProfileQualification, error) {
	var profileQualification models.LabourProfileQualification
	err := r.db.WithContext(ctx).Preload("Qualification.Sport").Where("labour_profile_id = ? AND qualification_id = ?", labourProfileID, qualificationID).First(&profileQualification).Error
	if err != nil {
		return nil, err
	}
	return &profileQualification, nil
}

func (r *labourProfileQualificationRepository) Update(ctx context.Context, profileQualification *models.LabourProfileQualification) error {
	return r.db.WithContext(ctx).Save(profileQualification).Error
}

func (r *labourProfileQualificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.LabourProfileQualification{}).Error
}

func (r *labourProfileQualificationRepository) DeleteByLabourProfileID(ctx context.Context, labourProfileID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("labour_profile_id = ?", labourProfileID).Delete(&models.LabourProfileQualification{}).Error
}

func (r *labourProfileQualificationRepository) DeleteByLabourProfileAndQualification(ctx context.Context, labourProfileID, qualificationID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("labour_profile_id = ? AND qualification_id = ?", labourProfileID, qualificationID).Delete(&models.LabourProfileQualification{}).Error
}
