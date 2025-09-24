package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yakka-backend/internal/features/masters/payment_constants/entity/database"
	"github.com/yakka-backend/internal/features/masters/payment_constants/models"
	"github.com/yakka-backend/internal/features/masters/payment_constants/payload"
)

// PaymentConstantUsecase defines the interface for payment constant business logic
type PaymentConstantUsecase interface {
	CreateConstant(ctx context.Context, req payload.CreatePaymentConstantRequest) (*models.PaymentConstant, error)
	GetConstantByID(ctx context.Context, id uuid.UUID) (*models.PaymentConstant, error)
	GetConstantByName(ctx context.Context, name string) (*models.PaymentConstant, error)
	GetAllConstants(ctx context.Context) ([]*models.PaymentConstant, error)
	GetActiveConstants(ctx context.Context) ([]*models.PaymentConstant, error)
	UpdateConstant(ctx context.Context, id uuid.UUID, req payload.UpdatePaymentConstantRequest) (*models.PaymentConstant, error)
	UpdateConstantValue(ctx context.Context, name string, req payload.UpdatePaymentConstantValueRequest) (*models.PaymentConstant, error)
	DeleteConstant(ctx context.Context, id uuid.UUID) error
}

// paymentConstantUsecase implements PaymentConstantUsecase
type paymentConstantUsecase struct {
	paymentConstantRepo database.PaymentConstantRepository
}

// NewPaymentConstantUsecase creates a new payment constant usecase
func NewPaymentConstantUsecase(paymentConstantRepo database.PaymentConstantRepository) PaymentConstantUsecase {
	return &paymentConstantUsecase{
		paymentConstantRepo: paymentConstantRepo,
	}
}

// CreateConstant creates a new payment constant
func (u *paymentConstantUsecase) CreateConstant(ctx context.Context, req payload.CreatePaymentConstantRequest) (*models.PaymentConstant, error) {
	// Check if constant with same name already exists
	existingConstant, err := u.paymentConstantRepo.GetByName(ctx, req.Name)
	if err == nil && existingConstant != nil {
		return nil, fmt.Errorf("payment constant with name '%s' already exists", req.Name)
	}

	// Set default values
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	constant := &models.PaymentConstant{
		Name:        req.Name,
		Value:       req.Value,
		Description: req.Description,
		IsActive:    isActive,
	}

	if err := u.paymentConstantRepo.Create(ctx, constant); err != nil {
		return nil, fmt.Errorf("failed to create payment constant: %w", err)
	}

	return constant, nil
}

// GetConstantByID retrieves a payment constant by ID
func (u *paymentConstantUsecase) GetConstantByID(ctx context.Context, id uuid.UUID) (*models.PaymentConstant, error) {
	constant, err := u.paymentConstantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment constant: %w", err)
	}
	return constant, nil
}

// GetConstantByName retrieves a payment constant by name
func (u *paymentConstantUsecase) GetConstantByName(ctx context.Context, name string) (*models.PaymentConstant, error) {
	constant, err := u.paymentConstantRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment constant: %w", err)
	}
	return constant, nil
}

// GetAllConstants retrieves all payment constants
func (u *paymentConstantUsecase) GetAllConstants(ctx context.Context) ([]*models.PaymentConstant, error) {
	constants, err := u.paymentConstantRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment constants: %w", err)
	}
	return constants, nil
}

// GetActiveConstants retrieves all active payment constants
func (u *paymentConstantUsecase) GetActiveConstants(ctx context.Context) ([]*models.PaymentConstant, error) {
	constants, err := u.paymentConstantRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active payment constants: %w", err)
	}
	return constants, nil
}

// UpdateConstant updates a payment constant
func (u *paymentConstantUsecase) UpdateConstant(ctx context.Context, id uuid.UUID, req payload.UpdatePaymentConstantRequest) (*models.PaymentConstant, error) {
	// Get existing constant
	constant, err := u.paymentConstantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("payment constant not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if new name already exists (excluding current constant)
		existingConstant, err := u.paymentConstantRepo.GetByName(ctx, *req.Name)
		if err == nil && existingConstant != nil && existingConstant.ID != constant.ID {
			return nil, fmt.Errorf("payment constant with name '%s' already exists", *req.Name)
		}
		constant.Name = *req.Name
	}

	if req.Value != nil {
		constant.Value = *req.Value
	}

	if req.Description != nil {
		constant.Description = req.Description
	}

	if req.IsActive != nil {
		constant.IsActive = *req.IsActive
	}

	if err := u.paymentConstantRepo.Update(ctx, constant); err != nil {
		return nil, fmt.Errorf("failed to update payment constant: %w", err)
	}

	return constant, nil
}

// UpdateConstantValue updates only the value of a payment constant by name
func (u *paymentConstantUsecase) UpdateConstantValue(ctx context.Context, name string, req payload.UpdatePaymentConstantValueRequest) (*models.PaymentConstant, error) {
	// Get existing constant
	constant, err := u.paymentConstantRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("payment constant not found: %w", err)
	}

	// Update value
	constant.Value = req.Value

	if err := u.paymentConstantRepo.Update(ctx, constant); err != nil {
		return nil, fmt.Errorf("failed to update payment constant value: %w", err)
	}

	return constant, nil
}

// DeleteConstant deletes a payment constant
func (u *paymentConstantUsecase) DeleteConstant(ctx context.Context, id uuid.UUID) error {
	// Check if constant exists
	_, err := u.paymentConstantRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("payment constant not found: %w", err)
	}

	if err := u.paymentConstantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete payment constant: %w", err)
	}

	return nil
}
