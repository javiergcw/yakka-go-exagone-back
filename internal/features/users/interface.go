package users

import (
	"context"

	"github.com/yakka-backend/internal/features/users/models"
)

// Repository defines the interface for user data operations
type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// UseCase defines the interface for user business logic
type UseCase interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id uint) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
}
