package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/payment_constants/models"
)

// PaymentConstantRepository defines the interface for payment constant operations
type PaymentConstantRepository interface {
	Create(ctx context.Context, constant *models.PaymentConstant) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PaymentConstant, error)
	GetByName(ctx context.Context, name string) (*models.PaymentConstant, error)
	GetAll(ctx context.Context) ([]*models.PaymentConstant, error)
	GetActive(ctx context.Context) ([]*models.PaymentConstant, error)
	Update(ctx context.Context, constant *models.PaymentConstant) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateValue(ctx context.Context, name string, value int) error
}
