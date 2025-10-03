package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/builder_profiles/models"
	"gorm.io/gorm"
)

// CompanyRepository defines the interface for company database operations
type CompanyRepository interface {
	Create(ctx context.Context, company *models.Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error)
	GetByName(ctx context.Context, name string) (*models.Company, error)
	GetAll(ctx context.Context) ([]*models.Company, error)
	Update(ctx context.Context, company *models.Company) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// companyRepository implements CompanyRepository
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

// Create creates a new company
func (r *companyRepository) Create(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

// GetByID retrieves a company by ID
func (r *companyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	var company models.Company
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// GetByName retrieves a company by name
func (r *companyRepository) GetByName(ctx context.Context, name string) (*models.Company, error) {
	var company models.Company
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// Update updates a company
func (r *companyRepository) Update(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Save(company).Error
}

// GetAll retrieves all companies
func (r *companyRepository) GetAll(ctx context.Context) ([]*models.Company, error) {
	var companies []*models.Company
	err := r.db.WithContext(ctx).Find(&companies).Error
	if err != nil {
		return nil, err
	}
	return companies, nil
}

// Delete deletes a company by ID
func (r *companyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Company{}).Error
}
