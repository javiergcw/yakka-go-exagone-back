package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/payment_constants/models"
	"gorm.io/gorm"
)

// paymentConstantRepository implements PaymentConstantRepository
type paymentConstantRepository struct {
	db *gorm.DB
}

// NewPaymentConstantRepository creates a new payment constant repository
func NewPaymentConstantRepository(db *gorm.DB) PaymentConstantRepository {
	return &paymentConstantRepository{
		db: db,
	}
}

// Create creates a new payment constant
func (r *paymentConstantRepository) Create(ctx context.Context, constant *models.PaymentConstant) error {
	return r.db.WithContext(ctx).Create(constant).Error
}

// GetByID retrieves a payment constant by ID
func (r *paymentConstantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PaymentConstant, error) {
	var constant models.PaymentConstant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&constant).Error
	if err != nil {
		return nil, err
	}
	return &constant, nil
}

// GetByName retrieves a payment constant by name
func (r *paymentConstantRepository) GetByName(ctx context.Context, name string) (*models.PaymentConstant, error) {
	var constant models.PaymentConstant
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&constant).Error
	if err != nil {
		return nil, err
	}
	return &constant, nil
}

// GetAll retrieves all payment constants
func (r *paymentConstantRepository) GetAll(ctx context.Context) ([]*models.PaymentConstant, error) {
	var constants []*models.PaymentConstant
	err := r.db.WithContext(ctx).Find(&constants).Error
	if err != nil {
		return nil, err
	}
	return constants, nil
}

// GetActive retrieves all active payment constants
func (r *paymentConstantRepository) GetActive(ctx context.Context) ([]*models.PaymentConstant, error) {
	var constants []*models.PaymentConstant
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&constants).Error
	if err != nil {
		return nil, err
	}
	return constants, nil
}

// Update updates a payment constant
func (r *paymentConstantRepository) Update(ctx context.Context, constant *models.PaymentConstant) error {
	return r.db.WithContext(ctx).Save(constant).Error
}

// Delete deletes a payment constant by ID
func (r *paymentConstantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PaymentConstant{}).Error
}

// UpdateValue updates the value of a payment constant by name
func (r *paymentConstantRepository) UpdateValue(ctx context.Context, name string, value int) error {
	return r.db.WithContext(ctx).Model(&models.PaymentConstant{}).
		Where("name = ?", name).
		Update("value", value).Error
}
